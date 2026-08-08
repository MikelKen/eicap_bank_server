package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CashCount struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Type     string
	Quantity int
	Subtotal decimal.Decimal

	DenominationID uuid.UUID    `gorm:"type:uuid;not null;"`
	Denomination   Denomination `gorm:"foreignKey:DenominationID;references:ID"`

	CashSessionID uuid.UUID   `gorm:"type:uuid;not null;"`
	CashSession   CashSession `gorm:"foreignKey:CashSessionID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *CashCount) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (CashCount) TableName() string {
	return "cash_counts"
}
