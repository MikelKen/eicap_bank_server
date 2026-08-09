package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BankOperation struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Code            string
	Date            time.Time
	PreviousBalance decimal.Decimal
	Import          decimal.Decimal
	EndBalance      decimal.Decimal

	TypeOperationID uuid.UUID     `gorm:"type:uuid;"`
	TypeOperation   TypeOperation `gorm:"foreignKey:TypeOperationID;references:ID"`

	CashSessionID *uuid.UUID   `gorm:"type:uuid;"`
	CashSession   *CashSession `gorm:"foreignKey:CashSessionID;references:ID"`

	AccountID *uuid.UUID `gorm:"type:uuid;"`
	Account   *Account   `gorm:"foreignKey:AccountID;references:ID"`

	OperationInformation *OperationInformation `gorm:"foreignKey:BankOperationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *BankOperation) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (BankOperation) TableName() string {
	return "bank_operations"
}
