package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Denomination struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Type  string
	Value decimal.Decimal //200, 100, 50, 20, 10, 5, 2, 1
	Name  string

	CashCounts []CashCount `gorm:"foreignKey:DenominationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *Denomination) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (Denomination) TableName() string {
	return "denominations"
}
