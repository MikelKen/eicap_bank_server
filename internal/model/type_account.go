package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TypeAccount struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Name string

	Accounts []Account `gorm:"foreignKey:TypeAccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *TypeAccount) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (TypeAccount) TableName() string {
	return "type_accounts"
}
