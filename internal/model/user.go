package model

import (
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/enum"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Name     string
	UserName *string `gorm:"uniqueIndex"`
	Email    *string `gorm:"uniqueIndex"`
	Password string
	Avatar   string
	Role     enum.Permission `gorm:"type:enum('admin','student');default:'student';not null;"`

	Clients      []Client      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CashSessions []CashSession `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (User) TableName() string {
	return "users"
}
