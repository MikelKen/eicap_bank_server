package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperationInformation struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Origin      string
	Reason      string
	Destination string
	Details     string

	BankOperationID uuid.UUID      `gorm:"type:uuid;"`
	BankOperation   *BankOperation `gorm:"foreignKey:BankOperationID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *OperationInformation) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (OperationInformation) TableName() string {
	return "operation_informations"
}
