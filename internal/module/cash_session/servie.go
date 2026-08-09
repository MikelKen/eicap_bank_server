package cashsession

import (
	"context"
	"errors"
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/internal/module/denomination"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Service interface {
	Open(ctx context.Context, userID uuid.UUID, input *Open) error
	Close(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, input *Close) error
	FindByID(ctx context.Context, id uuid.UUID) (*response.CashSession, error)
	FindAll(ctx context.Context, filter CashSessionFilter) (pagination.Response[response.CashSession], error)
	FindOpenByUserID(ctx context.Context, userID uuid.UUID) (*response.CashSession, error)
}

type service struct {
	repo             Repo
	denominationRepo denomination.Repo
}

func NewService(repo Repo, denominationRepo denomination.Repo) Service {
	return &service{repo: repo, denominationRepo: denominationRepo}
}

func (s *service) Open(ctx context.Context, userID uuid.UUID, input *Open) error {
	// 1. Un usuario no puede abrir dos cajas al mismo tiempo.
	if _, err := s.repo.FindOpenByUserID(ctx, userID); err == nil {
		return response.Conflict("Ya tienes una sesión de caja abierta")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	counts, total, err := s.buildCounts(ctx, input.Counts, CountTypeOpening)
	if err != nil {
		return err
	}

	session := &model.CashSession{
		State:         StateOpen,
		OpeningDate:   time.Now(),
		OpeningAmount: total,
		UserID:        userID,
	}

	return s.repo.CreateWithCounts(ctx, session, counts)
}

func (s *service) Close(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, input *Close) error {
	session, err := s.repo.FindByID(ctx, sessionID)
	if err != nil {
		return response.NotFound("Sesión de caja no encontrada")
	}

	if session.UserID != userID {
		return response.Forbidden("No puedes cerrar una caja que no te pertenece")
	}

	if session.State != StateOpen {
		return response.Conflict("Esta sesión de caja ya está cerrada")
	}

	counts, total, err := s.buildCounts(ctx, input.Counts, CountTypeClosing)
	if err != nil {
		return err
	}

	// TODO: cuando exista el módulo de bank_operations, el "expected" debe calcularse como
	// OpeningAmount + depósitos - retiros registrados durante la sesión.
	// Por ahora se usa OpeningAmount como placeholder.
	expected := session.OpeningAmount

	session.State = StateClosed
	session.ClosingDate = time.Now()
	session.ClosingAmount = total
	session.ExpectedAmount = expected
	session.DifferenceAmount = total.Sub(expected)

	return s.repo.UpdateWithCounts(ctx, session, counts)
}

// buildCounts valida las denominaciones, calcula subtotales y arma los model.CashCount + total.
func (s *service) buildCounts(ctx context.Context, inputs []CountInput, typ string) ([]model.CashCount, decimal.Decimal, error) {
	seen := make(map[string]bool, len(inputs))
	counts := make([]model.CashCount, 0, len(inputs))
	total := decimal.Zero

	for _, in := range inputs {
		denominationID, err := uuid.Parse(in.DenominationID)
		if err != nil {
			return nil, decimal.Zero, response.BadRequest("ID de denominación inválido")
		}

		if seen[in.DenominationID] {
			return nil, decimal.Zero, response.BadRequest("No puedes repetir la misma denominación en el conteo")
		}
		seen[in.DenominationID] = true

		den, err := s.denominationRepo.FindByID(denominationID)
		if err != nil {
			return nil, decimal.Zero, response.NotFound("Una de las denominaciones no existe")
		}

		subtotal := den.Value.Mul(decimal.NewFromInt(int64(in.Quantity)))
		total = total.Add(subtotal)

		counts = append(counts, model.CashCount{
			Type:           typ,
			Quantity:       in.Quantity,
			Subtotal:       subtotal,
			DenominationID: denominationID,
		})
	}

	return counts, total, nil
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*response.CashSession, error) {
	session, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, response.NotFound("Sesión de caja no encontrada")
	}
	return response.CashSessionToResponse(session), nil
}

func (s *service) FindAll(ctx context.Context, filter CashSessionFilter) (pagination.Response[response.CashSession], error) {
	sessions, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return pagination.Response[response.CashSession]{}, err
	}

	items := response.CashSessionsToResponse(sessions)
	return pagination.NewResponse(items, total, filter.Params), nil
}

func (s *service) FindOpenByUserID(ctx context.Context, userID uuid.UUID) (*response.CashSession, error) {
	session, err := s.repo.FindOpenByUserID(ctx, userID)
	if err != nil {
		return nil, response.NotFound("No tienes una sesión de caja abierta")
	}
	return response.CashSessionToResponse(session), nil
}
