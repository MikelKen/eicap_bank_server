package bankoperation

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/internal/module/account"
	cashsession "github.com/Eicap/EICAP-BANK/server/internal/module/cash_session"
	"github.com/Eicap/EICAP-BANK/server/internal/module/client"
	typeoperation "github.com/Eicap/EICAP-BANK/server/internal/module/type_operation"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, input *Create) error
	// CreateAccountOpening registra la apertura de una cuenta (APC).
	CreateAccountOpening(ctx context.Context, accountID uuid.UUID) error
	// RecordCashOpening registra la apertura de caja (APCA).
	RecordCashOpening(ctx context.Context, sessionID uuid.UUID, amount decimal.Decimal) error
	// RecordCashClosing registra el cierre de caja (CICA).
	RecordCashClosing(ctx context.Context, sessionID uuid.UUID, amount decimal.Decimal) error
	// SessionTotals devuelve ingresos y egresos acumulados en una sesión de caja.
	SessionTotals(ctx context.Context, sessionID uuid.UUID) (decimal.Decimal, decimal.Decimal, error)
	FindByID(ctx context.Context, id uuid.UUID) (*response.BankOperation, error)
	FindAll(ctx context.Context, filter BankOperationFilter) (pagination.Response[response.BankOperation], error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID, filter BankOperationFilter) (pagination.Response[response.BankOperation], error)
	// FindByActiveSession devuelve todas las operaciones realizadas durante la sesión de caja activa del usuario.
	FindByActiveSession(ctx context.Context, userID uuid.UUID, filter BankOperationFilter) (pagination.Response[response.BankOperation], error)
	// FindAllByClientID devuelve todas las operaciones de las cuentas de un cliente.
	FindAllByClientID(ctx context.Context, clientID uuid.UUID, filter BankOperationFilter) (pagination.Response[response.BankOperation], error)
}

type service struct {
	repo              Repo
	accountRepo       account.Repo
	typeOperationRepo typeoperation.Repo
	cashSessionRepo   cashsession.Repo
	clientRepo        client.Repo
}

func NewService(repo Repo, accountRepo account.Repo, typeOperationRepo typeoperation.Repo, cashSessionRepo cashsession.Repo, clientRepo client.Repo) Service {
	return &service{
		repo:              repo,
		accountRepo:       accountRepo,
		typeOperationRepo: typeOperationRepo,
		cashSessionRepo:   cashSessionRepo,
		clientRepo:        clientRepo,
	}
}

// Create registra una operación bancaria. Solo las operaciones de tipo ING/EGR
// afectan el balance de la cuenta y registran una OperationInformation; el resto
// de operaciones (APC, APCA, CICA...) se registran sin información adicional.
func (s *service) Create(ctx context.Context, userID uuid.UUID, input *Create) error {
	log.Println("Datos de la transaccion: ", *input)
	if input.Amount.IsNegative() {
		return response.BadRequest("El monto no puede ser negativo")
	}

	typeOp, err := s.typeOperationRepo.FindByCode(input.TypeOperationCode)
	if err != nil {
		return response.NotFound("Tipo de operación no encontrado")
	}

	switch input.TypeOperationCode {
	case CodeIncome, CodeExpense:
		return s.createIncomeOrExpense(ctx, userID, input, typeOp)
	default:
		return s.createOther(ctx, input, typeOp)
	}
}

func (s *service) createIncomeOrExpense(ctx context.Context, userID uuid.UUID, input *Create, typeOp *model.TypeOperation) error {
	if input.Amount.LessThanOrEqual(decimal.Zero) {
		return response.BadRequest("El monto debe ser mayor a cero")
	}

	if input.AccountID == nil {
		return response.BadRequest("La cuenta es obligatoria para esta operación")
	}

	if input.Info == nil {
		return response.BadRequest("La información de la operación es obligatoria")
	}
	if input.Info.Origin == "" || input.Info.Reason == "" || input.Info.Destination == "" {
		return response.BadRequest("La información de la operación está incompleta (origin, reason, destination)")
	}

	session, err := s.cashSessionRepo.FindOpenByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Conflict("Debes tener una caja abierta para registrar esta operación")
		}
		return err
	}

	acc, err := s.accountRepo.FindByID(*input.AccountID)
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
		CashSessionID:   &session.ID,
		AccountID:       input.AccountID,
	}

	info := &model.OperationInformation{
		Origin:      input.Info.Origin,
		Reason:      input.Info.Reason,
		Destination: input.Info.Destination,
		Details:     input.Info.Details,
	}

	return s.repo.Create(ctx, operation, info)
}

func (s *service) createOther(ctx context.Context, input *Create, typeOp *model.TypeOperation) error {
	var acc *model.Account
	if input.AccountID != nil {
		var err error
		acc, err = s.accountRepo.FindByID(*input.AccountID)
		if err != nil {
			return response.NotFound("Cuenta no encontrada")
		}
	}

	previousBalance := decimal.Zero
	if acc != nil {
		previousBalance = acc.Balance
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
		EndBalance:      previousBalance,
		TypeOperationID: typeOp.ID,
		AccountID:       input.AccountID,
	}

	// Las operaciones que no son ING/EGR no registran OperationInformation.
	return s.repo.Create(ctx, operation, nil)
}

func (s *service) FindAllByUserID(ctx context.Context, userID uuid.UUID, filter BankOperationFilter) (pagination.Response[response.BankOperation], error) {
	list, total, err := s.repo.FindAllByUserID(ctx, userID, filter)
	if err != nil {
		return pagination.Response[response.BankOperation]{}, err
	}
	return pagination.NewResponse(response.BankOperationsToResponse(list), total, filter.Params), nil
}

func (s *service) CreateAccountOpening(ctx context.Context, accountID uuid.UUID) error {
	typeOp, err := s.typeOperationRepo.FindByCode(CodeAccountOpen)
	if err != nil {
		return response.NotFound("Tipo de operación no encontrado")
	}

	code, err := s.repo.NextCode(ctx)
	if err != nil {
		return err
	}

	operation := &model.BankOperation{
		Code:            code,
		Date:            time.Now(),
		PreviousBalance: decimal.Zero,
		Import:          decimal.Zero,
		EndBalance:      decimal.Zero,
		TypeOperationID: typeOp.ID,
		AccountID:       &accountID,
	}

	return s.repo.Create(ctx, operation, nil)
}

func (s *service) RecordCashOpening(ctx context.Context, sessionID uuid.UUID, amount decimal.Decimal) error {
	return s.recordCash(ctx, sessionID, amount, CodeCashOpen)
}

func (s *service) RecordCashClosing(ctx context.Context, sessionID uuid.UUID, amount decimal.Decimal) error {
	return s.recordCash(ctx, sessionID, amount, CodeCashClose)
}

func (s *service) recordCash(ctx context.Context, sessionID uuid.UUID, amount decimal.Decimal, code string) error {
	typeOp, err := s.typeOperationRepo.FindByCode(code)
	if err != nil {
		return response.NotFound("Tipo de operación no encontrado")
	}

	operationCode, err := s.repo.NextCode(ctx)
	if err != nil {
		return err
	}

	operation := &model.BankOperation{
		Code:            operationCode,
		Date:            time.Now(),
		PreviousBalance: decimal.Zero,
		Import:          amount,
		EndBalance:      decimal.Zero,
		TypeOperationID: typeOp.ID,
		CashSessionID:   &sessionID,
	}

	return s.repo.Create(ctx, operation, nil)
}

func (s *service) SessionTotals(ctx context.Context, sessionID uuid.UUID) (decimal.Decimal, decimal.Decimal, error) {
	return s.repo.SessionTotals(ctx, sessionID)
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

// FindByActiveSession devuelve todas las operaciones registradas durante la sesión
// de caja que el usuario tiene abierta en este momento.
func (s *service) FindByActiveSession(ctx context.Context, userID uuid.UUID, filter BankOperationFilter) (pagination.Response[response.BankOperation], error) {
	session, err := s.cashSessionRepo.FindOpenByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pagination.Response[response.BankOperation]{}, response.NotFound("No tienes una sesión de caja abierta")
		}
		return pagination.Response[response.BankOperation]{}, err
	}

	list, total, err := s.repo.FindBySessionID(ctx, session.ID, filter)
	if err != nil {
		return pagination.Response[response.BankOperation]{}, err
	}
	return pagination.NewResponse(response.BankOperationsToResponse(list), total, filter.Params), nil
}

// FindAllByClientID devuelve todas las operaciones asociadas a las cuentas de un cliente.
func (s *service) FindAllByClientID(ctx context.Context, clientID uuid.UUID, filter BankOperationFilter) (pagination.Response[response.BankOperation], error) {
	if err := s.clientRepo.Exist(clientID); err != nil {
		return pagination.Response[response.BankOperation]{}, response.NotFound("Cliente no encontrado")
	}

	list, total, err := s.repo.FindAllByClientID(ctx, clientID, filter)
	if err != nil {
		return pagination.Response[response.BankOperation]{}, err
	}
	return pagination.NewResponse(response.BankOperationsToResponse(list), total, filter.Params), nil
}
