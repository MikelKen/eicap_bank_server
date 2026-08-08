package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Name      string
	Ci        string `gorm:"uniqueIndex"`
	Sex       string
	BirthDate string

	UserID uuid.UUID `gorm:"type:uuid;not null;"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	Accounts []Account `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *Client) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (Client) TableName() string {
	return "clients"
}
