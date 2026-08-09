package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Number   string
	Interest decimal.Decimal
	Balance  decimal.Decimal
	Status   string

	ClientID uuid.UUID `gorm:"type:uuid;"`
	Client   Client    `gorm:"foreignKey:ClientID;references:ID"`

	TypeAccountID uuid.UUID   `gorm:"type:uuid;"`
	TypeAccount   TypeAccount `gorm:"foreignKey:TypeAccountID;references:ID"`

	BankOperations []BankOperation `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *Account) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (Account) TableName() string {
	return "accounts"
}
