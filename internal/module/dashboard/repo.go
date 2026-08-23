package dashboard

import (
	"context"
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// UserRoleStat es el conteo de usuarios por rol.
type UserRoleStat struct {
	Role             string
	Total            int64
	CreatedThisMonth int64
}

type dayRow struct {
	Operations int64
	Income     decimal.Decimal
	Expense    decimal.Decimal
}

type Repo interface {
	// ClientsTotals devuelve el total de clientes y los creados en el mes actual.
	ClientsTotals(ctx context.Context, userID *uuid.UUID) (int64, int64, error)
	// OperationsTotals devuelve cantidad de operaciones y suma de ingresos/egresos.
	OperationsTotals(ctx context.Context, userID *uuid.UUID) (int64, decimal.Decimal, decimal.Decimal, error)
	// OperationsByType devuelve conteo y volumen por tipo (incluye tipos sin operaciones).
	OperationsByType(ctx context.Context, userID *uuid.UUID) ([]TypeCount, error)
	// CashSessionsStats devuelve aperturas, cierres y sesiones abiertas actualmente.
	CashSessionsStats(ctx context.Context, userID *uuid.UUID) (CashSessionsStats, error)
	// ActivityPerDay devuelve la actividad agregada por día desde `from`, indexada por YYYY-MM-DD.
	ActivityPerDay(ctx context.Context, from time.Time, userID *uuid.UUID) (map[string]dayRow, error)
	// RecentOperations devuelve las últimas operaciones registradas.
	RecentOperations(ctx context.Context, limit int, userID *uuid.UUID) ([]model.BankOperation, error)
	// UsersStats devuelve el conteo de usuarios por rol con altas del mes actual.
	UsersStats(ctx context.Context) ([]UserRoleStat, error)
	// UsersPerformance devuelve clientes, sesiones y operaciones por usuario.
	UsersPerformance(ctx context.Context, limit int) ([]UserPerformance, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) ClientsTotals(ctx context.Context, userID *uuid.UUID) (int64, int64, error) {
	var row struct {
		Total        int64
		NewThisMonth int64
	}

	q := r.db.WithContext(ctx).Model(&model.Client{}).
		Select(`COUNT(*) AS total,
			COUNT(*) FILTER (WHERE clients.created_at >= DATE_TRUNC('month', NOW())) AS new_this_month`)
	if userID != nil {
		q = q.Where("clients.user_id = ?", *userID)
	}

	if err := q.Scan(&row).Error; err != nil {
		return 0, 0, err
	}
	return row.Total, row.NewThisMonth, nil
}

// operationScope filtra operaciones por usuario: las registradas en sus sesiones
// de caja más las asociadas a cuentas de sus propios clientes (ej. aperturas).
func operationScope(db *gorm.DB, userID *uuid.UUID) *gorm.DB {
	db = db.Model(&model.BankOperation{})
	if userID != nil {
		return db.Where(
			"bank_operations.cash_session_id IN (?) OR bank_operations.account_id IN (?)",
			db.Session(&gorm.Session{NewDB: true}).Model(&model.CashSession{}).
				Select("id").Where("user_id = ?", *userID),
			db.Session(&gorm.Session{NewDB: true}).Model(&model.Account{}).
				Select("id").Where("client_id IN (?)",
					db.Session(&gorm.Session{NewDB: true}).Model(&model.Client{}).
						Select("id").Where("user_id = ?", *userID),
				),
		)
	}
	return db
}

func (r *repo) OperationsTotals(ctx context.Context, userID *uuid.UUID) (int64, decimal.Decimal, decimal.Decimal, error) {
	var row struct {
		Total   int64
		Income  decimal.Decimal
		Expense decimal.Decimal
	}

	err := operationScope(r.db.WithContext(ctx).Session(&gorm.Session{}), userID).
		Select(`COUNT(*) AS total,
			COALESCE(SUM(bank_operations.import) FILTER (WHERE type_operation.code = 'ING'), 0) AS income,
			COALESCE(SUM(bank_operations.import) FILTER (WHERE type_operation.code = 'EGR'), 0) AS expense`).
		Joins("LEFT JOIN type_operations type_operation ON type_operation.id = bank_operations.type_operation_id AND type_operation.deleted_at IS NULL").
		Scan(&row).Error

	return row.Total, row.Income, row.Expense, err
}

func (r *repo) OperationsByType(ctx context.Context, userID *uuid.UUID) ([]TypeCount, error) {
	// Un solo LEFT JOIN desde type_operations para incluir tipos sin operaciones con 0.
	// La condición de alcance va dentro del ON para no descartar esos tipos.
	join := `LEFT JOIN bank_operations ON bank_operations.type_operation_id = type_operations.id
		AND bank_operations.deleted_at IS NULL`
	args := []any{}
	if userID != nil {
		join += ` AND (bank_operations.cash_session_id IN (SELECT id FROM cash_sessions WHERE deleted_at IS NULL AND user_id = ?)
			OR bank_operations.account_id IN (SELECT id FROM accounts WHERE deleted_at IS NULL AND client_id IN (SELECT id FROM clients WHERE deleted_at IS NULL AND user_id = ?)))`
		args = append(args, *userID, *userID)
	}

	var rows []struct {
		Code  string
		Name  string
		Count int64
		Total decimal.Decimal
	}

	err := r.db.WithContext(ctx).Session(&gorm.Session{NewDB: true}).
		Table("type_operations").
		Select(`type_operations.code AS code,
			type_operations.name AS name,
			COALESCE(COUNT(bank_operations.id), 0) AS count,
			COALESCE(SUM(bank_operations.import), 0) AS total`).
		Joins(join, args...).
		Where("type_operations.deleted_at IS NULL").
		Group("type_operations.code, type_operations.name").
		Order("type_operations.code ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]TypeCount, 0, len(rows))
	for _, row := range rows {
		result = append(result, TypeCount{
			Code:  row.Code,
			Name:  row.Name,
			Count: row.Count,
			Total: row.Total.StringFixed(2),
		})
	}
	return result, nil
}

func (r *repo) CashSessionsStats(ctx context.Context, userID *uuid.UUID) (CashSessionsStats, error) {
	var row struct {
		Openings      int64
		Closings      int64
		CurrentlyOpen int64
	}

	q := r.db.WithContext(ctx).Model(&model.CashSession{}).
		Select(`COUNT(*) AS openings,
			COUNT(*) FILTER (WHERE cash_sessions.state = 'closed') AS closings,
			COUNT(*) FILTER (WHERE cash_sessions.state = 'open') AS currently_open`)
	if userID != nil {
		q = q.Where("cash_sessions.user_id = ?", *userID)
	}

	if err := q.Scan(&row).Error; err != nil {
		return CashSessionsStats{}, err
	}
	return CashSessionsStats{
		Openings:      row.Openings,
		Closings:      row.Closings,
		CurrentlyOpen: row.CurrentlyOpen,
	}, nil
}

// ActivityPerDay devuelve la actividad agregada por día desde `from`, indexada por YYYY-MM-DD.
func (r *repo) ActivityPerDay(ctx context.Context, from time.Time, userID *uuid.UUID) (map[string]dayRow, error) {
	var rows []struct {
		Day        time.Time
		Operations int64
		Income     decimal.Decimal
		Expense    decimal.Decimal
	}

	err := operationScope(r.db.WithContext(ctx).Session(&gorm.Session{}), userID).
		Select(`DATE(bank_operations.date) AS day,
			COUNT(*) AS operations,
			COALESCE(SUM(bank_operations.import) FILTER (WHERE type_operation.code = 'ING'), 0) AS income,
			COALESCE(SUM(bank_operations.import) FILTER (WHERE type_operation.code = 'EGR'), 0) AS expense`).
		Joins("LEFT JOIN type_operations type_operation ON type_operation.id = bank_operations.type_operation_id AND type_operation.deleted_at IS NULL").
		Where("bank_operations.date >= ?", from).
		Group("day").
		Order("day ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]dayRow, len(rows))
	for _, row := range rows {
		result[row.Day.Format("2006-01-02")] = dayRow{
			Operations: row.Operations,
			Income:     row.Income,
			Expense:    row.Expense,
		}
	}
	return result, nil
}

func (r *repo) RecentOperations(ctx context.Context, limit int, userID *uuid.UUID) ([]model.BankOperation, error) {
	var ops []model.BankOperation
	err := operationScope(r.db.WithContext(ctx).Session(&gorm.Session{}), userID).
		Preload("TypeOperation").
		Order("bank_operations.date DESC").
		Limit(limit).
		Find(&ops).Error
	return ops, err
}

func (r *repo) UsersStats(ctx context.Context) ([]UserRoleStat, error) {
	var rows []UserRoleStat
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Select(`role::text AS role,
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE users.created_at >= DATE_TRUNC('month', NOW())) AS created_this_month`).
		Group("role").
		Scan(&rows).Error
	return rows, err
}

func (r *repo) UsersPerformance(ctx context.Context, limit int) ([]UserPerformance, error) {
	var rows []UserPerformance
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT u.id::text AS id,
			       u.name AS name,
			       u.role::text AS role,
			       u.created_at AS created_at,
			       (SELECT COUNT(*) FROM clients c WHERE c.user_id = u.id AND c.deleted_at IS NULL) AS clients,
			       (SELECT COUNT(*) FROM cash_sessions s WHERE s.user_id = u.id AND s.deleted_at IS NULL) AS sessions,
			       (SELECT COUNT(*) FROM bank_operations bo
			         WHERE bo.deleted_at IS NULL
			           AND bo.cash_session_id IN (
			             SELECT s2.id FROM cash_sessions s2 WHERE s2.user_id = u.id AND s2.deleted_at IS NULL)) AS operations
			FROM users u
			WHERE u.deleted_at IS NULL
			ORDER BY clients DESC, operations DESC, u.created_at ASC
			LIMIT ?`,
			limit).
		Scan(&rows).Error
	return rows, err
}
