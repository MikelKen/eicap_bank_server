package cashsession

import (
	"context"
	"errors"

	"github.com/Eicap/EICAP-BANK/server/internal/generated"
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo interface {
	// CreateWithCounts crea la sesión y sus conteos de apertura en una sola transacción.
	CreateWithCounts(ctx context.Context, session *model.CashSession, counts []model.CashCount) error
	// UpdateWithCounts actualiza la sesión (cierre) e inserta los conteos de cierre en una transacción.
	UpdateWithCounts(ctx context.Context, session *model.CashSession, counts []model.CashCount) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.CashSession, error)
	FindOpenByUserID(ctx context.Context, userID uuid.UUID) (*model.CashSession, error)
	FindAll(ctx context.Context, filter CashSessionFilter) ([]model.CashSession, int64, error)
	Exist(id uuid.UUID) error
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return &repo{db: db}
}

func (r *repo) CreateWithCounts(ctx context.Context, session *model.CashSession, counts []model.CashCount) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(session).Error; err != nil {
			return err
		}

		for i := range counts {
			counts[i].CashSessionID = session.ID
		}

		if len(counts) > 0 {
			if err := tx.Create(&counts).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *repo) UpdateWithCounts(ctx context.Context, session *model.CashSession, counts []model.CashCount) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := range counts {
			counts[i].CashSessionID = session.ID
		}

		if len(counts) > 0 {
			if err := tx.Create(&counts).Error; err != nil {
				return err
			}
		}

		if err := tx.Save(session).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repo) FindByID(ctx context.Context, id uuid.UUID) (*model.CashSession, error) {
	var session model.CashSession
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("CashCounts.Denomination").
		Where(generated.CashSession.ID.Eq(id)).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repo) FindOpenByUserID(ctx context.Context, userID uuid.UUID) (*model.CashSession, error) {
	var session model.CashSession
	if err := r.db.WithContext(ctx).
		Where(generated.CashSession.UserID.Eq(userID)).
		Where(generated.CashSession.State.Eq(StateOpen)).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repo) FindAll(ctx context.Context, filter CashSessionFilter) ([]model.CashSession, int64, error) {
	return pagination.GormPaginate[model.CashSession](
		r.db.WithContext(ctx).Model(&model.CashSession{}).
			Preload("User"),
		filter.Params,
		func(db *gorm.DB) *gorm.DB {
			if filter.State != "" {
				db = db.Where("state = ?", filter.State)
			}

			order := "desc"
			if filter.Order == "asc" {
				order = "asc"
			}
			return db.Order("opening_date " + order)
		},
	)
}

func (r *repo) Exist(id uuid.UUID) error {
	var count int64
	if err := r.db.Model(&model.CashSession{}).
		Where(generated.CashSession.ID.Eq(id)).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return response.NotFound("Sesión de caja no encontrada")
	}
	return nil
}

var ErrNotFound = errors.New("cash session not found")
