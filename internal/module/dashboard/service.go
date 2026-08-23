package dashboard

import (
	"context"
	"sync"
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/enum"
	"github.com/google/uuid"
)

const (
	activityDays = 14
	recentOpsMax = 8
	topUsersMax  = 8
)

type Service interface {
	Summary(ctx context.Context, userID uuid.UUID, role enum.Permission) (*Summary, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

// Summary arma el resumen del panel. Los estudiantes ven solo sus propios datos;
// los administradores ven datos globales más estadísticas de usuarios.
// Las consultas se ejecutan en paralelo para minimizar la latencia total.
func (s *service) Summary(ctx context.Context, userID uuid.UUID, role enum.Permission) (*Summary, error) {
	var scope *uuid.UUID
	if role != enum.Admin {
		scope = &userID
	}

	summary := &Summary{
		Role:             role.String(),
		Clients:          ClientsStats{},
		Operations:       OperationsStats{ByType: []TypeCount{}},
		Activity:         []ActivityPoint{},
		RecentOperations: []RecentOperation{},
	}

	steps := []func(context.Context) error{
		s.fillClients(summary, scope),
		s.fillOperationTotals(summary, scope),
		s.fillOperationsByType(summary, scope),
		s.fillCashSessions(summary, scope),
		s.fillActivity(summary, scope),
		s.fillRecentOperations(summary, scope),
	}
	if role == enum.Admin {
		// Pre-inicializado para que cada paso escriba solo sus propios campos.
		summary.Admin = &AdminStats{UsersByRole: []NameCount{}, TopUsers: []UserPerformance{}}
		steps = append(steps,
			s.fillAdminUserStats(summary),
			s.fillAdminTopUsers(summary),
		)
	}

	errs := make([]error, len(steps))
	var wg sync.WaitGroup
	for i, step := range steps {
		wg.Add(1)
		go func(i int, step func(context.Context) error) {
			defer wg.Done()
			errs[i] = step(ctx)
		}(i, step)
	}
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return summary, nil
}

func (s *service) fillClients(summary *Summary, scope *uuid.UUID) func(context.Context) error {
	return func(ctx context.Context) error {
		total, newThisMonth, err := s.repo.ClientsTotals(ctx, scope)
		if err != nil {
			return err
		}
		summary.Clients = ClientsStats{Total: total, NewThisMonth: newThisMonth}
		return nil
	}
}

func (s *service) fillOperationTotals(summary *Summary, scope *uuid.UUID) func(context.Context) error {
	return func(ctx context.Context) error {
		total, income, expense, err := s.repo.OperationsTotals(ctx, scope)
		if err != nil {
			return err
		}
		summary.Operations.Total = total
		summary.Operations.IncomeTotal = income.StringFixed(2)
		summary.Operations.ExpenseTotal = expense.StringFixed(2)
		summary.Operations.NetTotal = income.Sub(expense).StringFixed(2)
		return nil
	}
}

func (s *service) fillOperationsByType(summary *Summary, scope *uuid.UUID) func(context.Context) error {
	return func(ctx context.Context) error {
		byType, err := s.repo.OperationsByType(ctx, scope)
		if err != nil {
			return err
		}
		summary.Operations.ByType = byType
		return nil
	}
}

func (s *service) fillCashSessions(summary *Summary, scope *uuid.UUID) func(context.Context) error {
	return func(ctx context.Context) error {
		stats, err := s.repo.CashSessionsStats(ctx, scope)
		if err != nil {
			return err
		}
		summary.CashSessions = stats
		return nil
	}
}

func (s *service) fillActivity(summary *Summary, scope *uuid.UUID) func(context.Context) error {
	return func(ctx context.Context) error {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		from := today.AddDate(0, 0, -(activityDays - 1))

		rows, err := s.repo.ActivityPerDay(ctx, from, scope)
		if err != nil {
			return err
		}

		points := make([]ActivityPoint, 0, activityDays)
		for i := activityDays - 1; i >= 0; i-- {
			day := today.AddDate(0, 0, -i)
			key := day.Format("2006-01-02")
			row := rows[key]
			points = append(points, ActivityPoint{
				Date:       key,
				Operations: row.Operations,
				Income:     row.Income.StringFixed(2),
				Expense:    row.Expense.StringFixed(2),
			})
		}
		summary.Activity = points
		return nil
	}
}

func (s *service) fillRecentOperations(summary *Summary, scope *uuid.UUID) func(context.Context) error {
	return func(ctx context.Context) error {
		ops, err := s.repo.RecentOperations(ctx, recentOpsMax, scope)
		if err != nil {
			return err
		}

		recent := make([]RecentOperation, 0, len(ops))
		for _, op := range ops {
			recent = append(recent, RecentOperation{
				ID:       op.ID.String(),
				Code:     op.Code,
				Date:     op.Date,
				TypeCode: op.TypeOperation.Code,
				TypeName: op.TypeOperation.Name,
				Import:   op.Import.StringFixed(2),
			})
		}
		summary.RecentOperations = recent
		return nil
	}
}

func (s *service) fillAdminUserStats(summary *Summary) func(context.Context) error {
	return func(ctx context.Context) error {
		stats, err := s.repo.UsersStats(ctx)
		if err != nil {
			return err
		}

		labels := map[string]string{"admin": "Administradores", "student": "Estudiantes"}
		for _, row := range stats {
			summary.Admin.UsersTotal += row.Total
			summary.Admin.UsersCreatedThisMonth += row.CreatedThisMonth
			label := labels[row.Role]
			if label == "" {
				label = row.Role
			}
			summary.Admin.UsersByRole = append(summary.Admin.UsersByRole, NameCount{Label: label, Value: row.Total})
		}
		return nil
	}
}

func (s *service) fillAdminTopUsers(summary *Summary) func(context.Context) error {
	return func(ctx context.Context) error {
		performance, err := s.repo.UsersPerformance(ctx, topUsersMax)
		if err != nil {
			return err
		}
		summary.Admin.TopUsers = performance
		return nil
	}
}
