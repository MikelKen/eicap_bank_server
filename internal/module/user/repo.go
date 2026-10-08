package user

import (
	"context"
	"strings"

	"github.com/Eicap/EICAP-BANK/server/internal/generated"
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo interface {
	Create(input *model.User) error
	Update(input *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByUserName(username string) (*model.User, error)
	FindByID(id uuid.UUID) (*model.User, error)
	FindAll(ctx context.Context, filter UserFilter) ([]model.User, int64, error)
	Exist(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) Create(input *model.User) error {
	return r.db.Create(input).Error
}

func (r *repo) Update(input *model.User) error {
	return r.db.Save(input).Error
}

func (r *repo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where(generated.User.Email.Eq(email)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) GetByUserName(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where(generated.User.UserName.Eq(username)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.Where(generated.User.ID.Eq(id)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) FindAll(ctx context.Context, filter UserFilter) ([]model.User, int64, error) {
	return pagination.GormPaginate[model.User](
		r.db.WithContext(ctx).Unscoped().Model(&model.User{}),
		filter.Params,
		func(db *gorm.DB) *gorm.DB {
			if filter.Name != "" {
				db = db.Where("name ILIKE ?", "%"+filter.Name+"%")
			}

			if filter.Roles != "" {
				roles := strings.Split(filter.Roles, ",")
				db = db.Where("role IN ?", roles)
			}
			if filter.Order != "asc" {
				filter.Order = "desc"
			}
			return db.Order("created_at " + filter.Order)
		},
	)
}

func (r *repo) Exist(id uuid.UUID) error {
	var count int64
	if err := r.db.Model(&model.User{}).
		Where(generated.User.ID.Eq(id)).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return response.NotFound("Usuario no encontrado")
	}
	return nil
}

func (r *repo) Delete(id uuid.UUID) error {
	return r.db.Where(generated.User.ID.Eq(id)).Delete(&model.User{}).Error
}
