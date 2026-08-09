package denomination

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
	Create(input *model.Denomination) error
	Update(input *model.Denomination) error
	FindByID(id uuid.UUID) (*model.Denomination, error)
	FindAll(ctx context.Context, filter DenominationFilter) ([]model.Denomination, int64, error)
	Exist(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) Create(input *model.Denomination) error {
	return r.db.Create(input).Error
}

func (r *repo) Update(input *model.Denomination) error {
	return r.db.Save(input).Error
}

func (r *repo) FindByID(id uuid.UUID) (*model.Denomination, error) {
	var denomination model.Denomination
	if err := r.db.Where(generated.Denomination.ID.Eq(id)).First(&denomination).Error; err != nil {
		return nil, err
	}
	return &denomination, nil
}

func (r *repo) FindAll(ctx context.Context, filter DenominationFilter) ([]model.Denomination, int64, error) {
	return pagination.GormPaginate[model.Denomination](
		r.db.WithContext(ctx).Model(&model.Denomination{}),
		filter.Params,
		func(db *gorm.DB) *gorm.DB {
			if filter.Type != "" {
				db = db.Where("type = ?", filter.Type)
			}

			order := "desc"
			if filter.Order == "asc" {
				order = "asc"
			}
			return db.Order("value " + order)
		},
	)
}

func (r *repo) Exist(id uuid.UUID) error {
	var count int64
	if err := r.db.Model(&model.Denomination{}).
		Where(generated.Denomination.ID.Eq(id)).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return response.NotFound("Denominación no encontrada")
	}
	return nil
}

func (r *repo) Delete(id uuid.UUID) error {
	return r.db.Where(generated.Denomination.ID.Eq(id)).Delete(&model.Denomination{}).Error
}
