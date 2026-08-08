package client

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/generated"
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo interface {
	Create(input *model.Client) error
	Update(input *model.Client) error
	GetClient(data string) (*model.Client, error) // get client by name o c.i.
	FindByID(id uuid.UUID) (*model.Client, error)
	FindAll(stx context.Context) ([]model.Client, int64, error)
	GetAllClientByUserID(userID uuid.UUID) (*model.Client, error)
	Exist(id uuid.UUID) error
	Delete(id uuid.UUID) error
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) Create(input *model.Client) error {
	return r.db.Create(input).Error
}

func (r *repo) Update(input *model.Client) error {
	return r.db.Save(input).Error
}

func (r *repo) GetClient(data string) (*model.Client, error) {
	var client model.Client
	if err := r.db.
		Where("name ILIKE ? OR ci = ?", "%"+data+"%", data).
		First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *repo) FindByID(id uuid.UUID) (*model.Client, error) {
	var client model.Client
	if err := r.db.Where(generated.Client.ID.Eq(id)).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *repo) FindAll(ctx context.Context) ([]model.Client, int64, error) {
	var clients []model.Client
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.Client{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&clients).Error; err != nil {
		return nil, 0, err
	}

	return clients, total, nil

}

func (r *repo) GetAllClientByUserID(userID uuid.UUID) (*model.Client, error) {
	var client model.Client
	if err := r.db.Where(generated.Client.UserID.Eq(userID)).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *repo) Exist(id uuid.UUID) error {
	var count int64
	if err := r.db.Model(&model.Client{}).
		Where(generated.Client.ID.Eq(id)).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return response.NotFound("Cliente no encontrado")
	}
	return nil
}

func (r *repo) Delete(id uuid.UUID) error {
	return r.db.Where(generated.Client.ID.Eq(id)).Delete(&model.Client{}).Error
}
