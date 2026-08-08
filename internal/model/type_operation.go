package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TypeOperation struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Code string
	Name string

	BankOperations []BankOperation `gorm:"foreignKey:TypeOperationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *TypeOperation) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (TypeOperation) TableName() string {
	return "type_operations"
}
