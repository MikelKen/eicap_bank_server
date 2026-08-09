package cashcount

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/generated"
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo interface {
	FindAllBySessionID(ctx context.Context, sessionID uuid.UUID) ([]model.CashCount, error)
	FindAllBySessionIDAndType(ctx context.Context, sessionID uuid.UUID, typ string) ([]model.CashCount, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) FindAllBySessionID(ctx context.Context, sessionID uuid.UUID) ([]model.CashCount, error) {
	var counts []model.CashCount
	if err := r.db.WithContext(ctx).
		Preload("Denomination").
		Where(generated.CashCount.CashSessionID.Eq(sessionID)).
		Order("created_at ASC").
		Find(&counts).Error; err != nil {
		return nil, err
	}
	return counts, nil
}

func (r *repo) FindAllBySessionIDAndType(ctx context.Context, sessionID uuid.UUID, typ string) ([]model.CashCount, error) {
	var counts []model.CashCount
	if err := r.db.WithContext(ctx).
		Preload("Denomination").
		Where(generated.CashCount.CashSessionID.Eq(sessionID)).
		Where(generated.CashCount.Type.Eq(typ)).
		Order("created_at ASC").
		Find(&counts).Error; err != nil {
		return nil, err
	}
	return counts, nil
}
