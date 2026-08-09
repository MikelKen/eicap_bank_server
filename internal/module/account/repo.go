package account

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Eicap/EICAP-BANK/server/internal/generated"
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo interface {
	Create(input *model.Account) error
	Update(input *model.Account) error
	FindByID(id uuid.UUID) (*model.Account, error)
	FindAll(ctx context.Context, filter AccountFilter) ([]model.Account, int64, error)
	FindAllByClientID(ctx context.Context, clientID uuid.UUID, filter AccountFilter) ([]model.Account, int64, error)
	NextNumber() (string, error)
	Exist(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) Create(input *model.Account) error {
	return r.db.Create(input).Error
}

func (r *repo) Update(input *model.Account) error {
	return r.db.Save(input).Error
}

func (r *repo) FindByID(id uuid.UUID) (*model.Account, error) {
	var acc model.Account
	if err := r.db.
		Preload("Client").
		Preload("TypeAccount").
		Where(generated.Account.ID.Eq(id)).
		First(&acc).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repo) FindAll(ctx context.Context, filter AccountFilter) ([]model.Account, int64, error) {
	return pagination.GormPaginate[model.Account](
		r.db.WithContext(ctx).Model(&model.Account{}).
			Preload("Client").Preload("TypeAccount"),
		filter.Params,
		func(db *gorm.DB) *gorm.DB {
			db = applyAccountFilters(db, filter)
			return db.Order("created_at " + orderDir(filter.Order))
		},
	)
}

func (r *repo) FindAllByClientID(ctx context.Context, clientID uuid.UUID, filter AccountFilter) ([]model.Account, int64, error) {
	return pagination.GormPaginate[model.Account](
		r.db.WithContext(ctx).Model(&model.Account{}).
			Preload("Client").Preload("TypeAccount").
			Where(generated.Account.ClientID.Eq(clientID)),
		filter.Params,
		func(db *gorm.DB) *gorm.DB {
			db = applyAccountFilters(db, filter)
			return db.Order("created_at " + orderDir(filter.Order))
		},
	)
}

func (r *repo) NextNumber() (string, error) {
	for {
		n := rand.Int63n(100_000_000)
		number := fmt.Sprintf("1%08d", n)

		var count int64
		if err := r.db.Unscoped().
			Model(&model.Account{}).
			Where("number = ?", number).
			Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return number, nil
		}
	}
}

func (r *repo) Exist(id uuid.UUID) error {
	var count int64
	if err := r.db.Model(&model.Account{}).
		Where(generated.Account.ID.Eq(id)).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return response.NotFound("Cuenta no encontrada")
	}
	return nil
}

func (r *repo) Delete(id uuid.UUID) error {
	return r.db.Where(generated.Account.ID.Eq(id)).Delete(&model.Account{}).Error
}

func applyAccountFilters(db *gorm.DB, filter AccountFilter) *gorm.DB {
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	if filter.ClientID != "" {
		db = db.Where("client_id = ?", filter.ClientID)
	}
	if filter.TypeAccountID != "" {
		db = db.Where("type_account_id = ?", filter.TypeAccountID)
	}
	return db
}

func orderDir(order string) string {
	if order == "asc" {
		return "asc"
	}
	return "desc"
}
