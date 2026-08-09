package bankoperation

import (
	"context"
	"errors"
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/internal/module/account"
	cashsession "github.com/Eicap/EICAP-BANK/server/internal/module/cash_session"
	typeoperation "github.com/Eicap/EICAP-BANK/server/internal/module/type_operation"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, input *Create) error
	FindByID(ctx context.Context, id uuid.UUID) (*response.BankOperation, error)
	FindAll(ctx context.Context, filter BankOperationFilter) (pagination.Response[response.BankOperation], error)
}

type service struct {
	repo              Repo
	accountRepo       account.Repo
	typeOperationRepo typeoperation.Repo
	cashSessionRepo   cashsession.Repo
}

func NewService(repo Repo, accountRepo account.Repo, typeOperationRepo typeoperation.Repo, cashSessionRepo cashsession.Repo) Service {
	return &service{
		repo:              repo,
		accountRepo:       accountRepo,
		typeOperationRepo: typeOperationRepo,
		cashSessionRepo:   cashSessionRepo,
	}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, input *Create) error {
	if input.Amount.LessThanOrEqual(decimal.Zero) {
		return response.BadRequest("El monto debe ser mayor a cero")
	}

	typeOp, err := s.typeOperationRepo.FindByCode(input.TypeOperationCode)
	if err != nil {
		return response.NotFound("Tipo de operación no encontrado")
	}

	session, err := s.cashSessionRepo.FindOpenByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Conflict("Debes tener una caja abierta para registrar esta operación")
		}
		return err
	}

	acc, err := s.accountRepo.FindByID(input.AccountID)
	if err != nil {
		return response.NotFound("Cuenta no encontrada")
	}

	previousBalance := acc.Balance
	var endBalance decimal.Decimal

	switch input.TypeOperationCode {
	case CodeIncome:
		endBalance = previousBalance.Add(input.Amount)
	case CodeExpense:
		if input.Amount.GreaterThan(previousBalance) {
			return response.BadRequest("Fondos insuficientes en la cuenta")
		}
		endBalance = previousBalance.Sub(input.Amount)
	default:
		// La validación del DTO (oneof=ING EGR) ya evita llegar aquí,
		// pero se deja como defensa extra.
		return response.BadRequest("Este endpoint solo admite operaciones de tipo ING o EGR")
	}

	code, err := s.repo.NextCode(ctx)
	if err != nil {
		return err
	}

	operation := &model.BankOperation{
		Code:            code,
		Date:            time.Now(),
		PreviousBalance: previousBalance,
		Import:          input.Amount,
		EndBalance:      endBalance,
		TypeOperationID: typeOp.ID,
		CashSessionID:   session.ID,
		AccountID:       input.AccountID,
	}

	// Solo ING/EGR llevan OperationInformation — el DTO ya lo exige como required,
	// así que aquí siempre viene no-nulo para este endpoint.
	info := &model.OperationInformation{
		Origin:      input.Info.Origin,
		Reason:      input.Info.Reason,
		Destination: input.Info.Destination,
		Details:     input.Info.Details,
	}

	return s.repo.CreateIncomeOrExpense(ctx, operation, info, endBalance)
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*response.BankOperation, error) {
	op, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, response.NotFound("Operación no encontrada")
	}
	return response.BankOperationToResponse(op), nil
}

func (s *service) FindAll(ctx context.Context, filter BankOperationFilter) (pagination.Response[response.BankOperation], error) {
	list, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return pagination.Response[response.BankOperation]{}, err
	}
	return pagination.NewResponse(response.BankOperationsToResponse(list), total, filter.Params), nil
}
