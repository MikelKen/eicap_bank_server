package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CashSession struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;"`
	State            string
	OpeningDate      time.Time
	OpeningAmount    decimal.Decimal
	ClosingDate      time.Time
	ClosingAmount    decimal.Decimal
	ExpectedAmount   decimal.Decimal
	DifferenceAmount decimal.Decimal

	UserID uuid.UUID `gorm:"type:uuid;not null;"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	BankOperations []BankOperation `gorm:"foreignKey:CashSessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CashCounts     []CashCount     `gorm:"foreignKey:CashSessionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *CashSession) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (CashSession) TableName() string {
	return "cash_sessions"
}
