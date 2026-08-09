package typeaccount

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/generated"
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo interface {
	Create(input *model.TypeAccount) error
	Update(input *model.TypeAccount) error
	FindByID(id uuid.UUID) (*model.TypeAccount, error)
	FindAll(ctx context.Context, filter TypeAccountFilter) ([]model.TypeAccount, int64, error)
	Exist(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) Create(input *model.TypeAccount) error {
	return r.db.Create(input).Error
}

func (r *repo) Update(input *model.TypeAccount) error {
	return r.db.Save(input).Error
}

func (r *repo) FindByID(id uuid.UUID) (*model.TypeAccount, error) {
	var typeAccount model.TypeAccount
	if err := r.db.Where(generated.TypeAccount.ID.Eq(id)).First(&typeAccount).Error; err != nil {
		return nil, err
	}
	return &typeAccount, nil
}

func (r *repo) FindAll(ctx context.Context, filter TypeAccountFilter) ([]model.TypeAccount, int64, error) {
	return pagination.GormPaginate[model.TypeAccount](
		r.db.WithContext(ctx).Model(&model.TypeAccount{}),
		filter.Params,
		func(db *gorm.DB) *gorm.DB {
			if filter.Name != "" {
				db = db.Where("name ILIKE ?", "%"+filter.Name+"%")
			}

			order := "desc"
			if filter.Order == "asc" {
				order = "asc"
			}
			return db.Order("created_at " + order)
		},
	)
}

func (r *repo) Exist(id uuid.UUID) error {
	var count int64
	if err := r.db.Model(&model.TypeAccount{}).
		Where(generated.TypeAccount.ID.Eq(id)).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return response.NotFound("Tipo de cuenta no encontrado")
	}
	return nil
}

func (r *repo) Delete(id uuid.UUID) error {
	return r.db.Where(generated.TypeAccount.ID.Eq(id)).Delete(&model.TypeAccount{}).Error
}
