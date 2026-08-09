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
	// CreateIncomeOrExpense crea la operación, opcionalmente su OperationInformation,
	// y actualiza el balance de la cuenta — todo en una sola transacción.
	CreateIncomeOrExpense(ctx context.Context, operation *model.BankOperation, info *model.OperationInformation, newBalance decimal.Decimal) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.BankOperation, error)
	FindAll(ctx context.Context, filter BankOperationFilter) ([]model.BankOperation, int64, error)
	NextCode(ctx context.Context) (string, error)
	ExistCode(ctx context.Context, code string) (bool, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) CreateIncomeOrExpense(ctx context.Context, operation *model.BankOperation, info *model.OperationInformation, newBalance decimal.Decimal) error {
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

		if err := tx.Model(&model.Account{}).
			Where("id = ?", operation.AccountID).
			Update("balance", newBalance).Error; err != nil {
			return err
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

			order := "desc"
			if filter.Order == "asc" {
				order = "asc"
			}
			return db.Order("bank_operations.date " + order)
		},
	)
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
