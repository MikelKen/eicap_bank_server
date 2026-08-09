package bankoperation

import (
	"context"
	"math/rand"

	"github.com/Eicap/EICAP-BANK/server/internal/generated"
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Repo interface {
	// Create registra la operación en una transacción: crea la operación, si `info` no es nil
	// crea su OperationInformation (solo ING/EGR), y si la operación tiene cuenta actualiza su balance.
	Create(ctx context.Context, operation *model.BankOperation, info *model.OperationInformation) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.BankOperation, error)
	FindAll(ctx context.Context, filter BankOperationFilter) ([]model.BankOperation, int64, error)
	// FindBySessionID devuelve las operaciones registradas durante una sesión de caja.
	FindBySessionID(ctx context.Context, sessionID uuid.UUID, filter BankOperationFilter) ([]model.BankOperation, int64, error)
	// FindAllByClientID devuelve las operaciones de las cuentas de un cliente.
	FindAllByClientID(ctx context.Context, clientID uuid.UUID, filter BankOperationFilter) ([]model.BankOperation, int64, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID, filter BankOperationFilter) ([]model.BankOperation, int64, error)
	// SessionTotals devuelve la suma de ingresos y egresos registrados durante una sesión de caja.
	SessionTotals(ctx context.Context, sessionID uuid.UUID) (income, expense decimal.Decimal, err error)
	NextCode(ctx context.Context) (string, error)
	ExistCode(ctx context.Context, code string) (bool, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, operation *model.BankOperation, info *model.OperationInformation) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(operation).Error; err != nil {
			return err
		}

		if info != nil {
			info.BankOperationID = operation.ID
			if err := tx.Create(info).Error; err != nil {
				return err
			}
		}

		if operation.AccountID != nil {
			if err := tx.Model(&model.Account{}).
				Where("id = ?", *operation.AccountID).
				Update("balance", operation.EndBalance).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *repo) FindByID(ctx context.Context, id uuid.UUID) (*model.BankOperation, error) {
	var op model.BankOperation
	if err := r.db.WithContext(ctx).
		Preload("TypeOperation").
		Preload("Account").
		Preload("OperationInformation").
		Where(generated.BankOperation.ID.Eq(id)).
		First(&op).Error; err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *repo) FindAll(ctx context.Context, filter BankOperationFilter) ([]model.BankOperation, int64, error) {
	return pagination.GormPaginate[model.BankOperation](
		r.db.WithContext(ctx).Model(&model.BankOperation{}).
			Preload("TypeOperation").Preload("Account").Preload("OperationInformation"),
		filter.Params,
		func(db *gorm.DB) *gorm.DB {
			if filter.AccountID != "" {
				db = db.Where("account_id = ?", filter.AccountID)
			}
			if filter.TypeOperationCode != "" {
				db = db.Joins("JOIN type_operations ON type_operations.id = bank_operations.type_operation_id").
					Where("type_operations.code = ?", filter.TypeOperationCode)
			}
			return orderByDate(db, filter)
		},
	)
}

func orderByDate(db *gorm.DB, filter BankOperationFilter) *gorm.DB {
	order := "desc"
	if filter.Order == "asc" {
		order = "asc"
	}
	return db.Order("bank_operations.date " + order)
}

func (r *repo) FindBySessionID(ctx context.Context, sessionID uuid.UUID, filter BankOperationFilter) ([]model.BankOperation, int64, error) {
	return pagination.GormPaginate[model.BankOperation](
		r.db.WithContext(ctx).Model(&model.BankOperation{}).
			Preload("TypeOperation").Preload("Account").Preload("OperationInformation").
			Where("cash_session_id = ?", sessionID),
		filter.Params,
		func(db *gorm.DB) *gorm.DB {
			if filter.TypeOperationCode != "" {
				db = db.Joins("JOIN type_operations ON type_operations.id = bank_operations.type_operation_id").
					Where("type_operations.code = ?", filter.TypeOperationCode)
			}
			return orderByDate(db, filter)
		},
	)
}

func (r *repo) FindAllByClientID(ctx context.Context, clientID uuid.UUID, filter BankOperationFilter) ([]model.BankOperation, int64, error) {
	return pagination.GormPaginate[model.BankOperation](
		r.db.WithContext(ctx).Model(&model.BankOperation{}).
			Preload("TypeOperation").Preload("Account").Preload("OperationInformation").
			Joins("JOIN accounts ON accounts.id = bank_operations.account_id").
			Where("accounts.client_id = ?", clientID),
		filter.Params,
		func(db *gorm.DB) *gorm.DB {
			if filter.AccountID != "" {
				db = db.Where("bank_operations.account_id = ?", filter.AccountID)
			}
			if filter.TypeOperationCode != "" {
				db = db.Joins("JOIN type_operations ON type_operations.id = bank_operations.type_operation_id").
					Where("type_operations.code = ?", filter.TypeOperationCode)
			}
			return orderByDate(db, filter)
		},
	)
}

func (r *repo) FindAllByUserID(ctx context.Context, userID uuid.UUID, filter BankOperationFilter) ([]model.BankOperation, int64, error) {
	return pagination.GormPaginate[model.BankOperation](
		r.db.WithContext(ctx).Model(&model.BankOperation{}).
			Preload("TypeOperation").Preload("Account").Preload("OperationInformation").
			Joins("JOIN cash_sessions ON cash_sessions.id = bank_operations.cash_session_id").
			Where("cash_sessions.user_id = ?", userID),
		filter.Params,
		func(db *gorm.DB) *gorm.DB {
			if filter.AccountID != "" {
				db = db.Where("bank_operations.account_id = ?", filter.AccountID)
			}
			if filter.TypeOperationCode != "" {
				db = db.Joins("JOIN type_operations ON type_operations.id = bank_operations.type_operation_id").
					Where("type_operations.code = ?", filter.TypeOperationCode)
			}
			return orderByDate(db, filter)
		},
	)
}

func (r *repo) SessionTotals(ctx context.Context, sessionID uuid.UUID) (income, expense decimal.Decimal, err error) {
	income = decimal.Zero
	expense = decimal.Zero

	type totalRow struct {
		Code  string
		Total decimal.Decimal
	}

	var rows []totalRow
	if err := r.db.WithContext(ctx).
		Model(&model.BankOperation{}).
		Select("type_operations.code AS code, COALESCE(SUM(bank_operations.import), 0) AS total").
		Joins("JOIN type_operations ON type_operations.id = bank_operations.type_operation_id").
		Where("bank_operations.cash_session_id = ?", sessionID).
		Where("type_operations.code IN ?", []string{CodeIncome, CodeExpense}).
		Group("type_operations.code").
		Scan(&rows).Error; err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	for _, row := range rows {
		switch row.Code {
		case CodeIncome:
			income = row.Total
		case CodeExpense:
			expense = row.Total
		}
	}

	return income, expense, nil
}

const codeLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const codeDigits = "0123456789"

// generateCode produce un código de 5 caracteres: 4 letras + 1 número (ej. "QWRT7").
func generateCode() string {
	code := make([]byte, 5)
	for i := 0; i < 4; i++ {
		code[i] = codeLetters[rand.Intn(len(codeLetters))]
	}
	code[4] = codeDigits[rand.Intn(len(codeDigits))]
	return string(code)
}

func (r *repo) ExistCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.BankOperation{}).
		Where(generated.BankOperation.Code.Eq(code)).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// NextCode genera un código único de 5 caracteres (4 letras + 1 número),
// reintentando en el raro caso de colisión.
func (r *repo) NextCode(ctx context.Context) (string, error) {
	for range 20 {
		code := generateCode()

		exists, err := r.ExistCode(ctx, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", gorm.ErrInvalidData
}
