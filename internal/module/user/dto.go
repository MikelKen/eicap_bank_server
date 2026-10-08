package user

import (
	"github.com/Eicap/EICAP-BANK/server/internal/enum"
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
)

type UserFilter struct {
	pagination.Params
	Name  string `form:"name" json:"name" query:"name"`
	Roles string `form:"roles" json:"roles" query:"roles"`
}

func (f *UserFilter) SetDefaults() {
	f.Params.SetDefaults()
}

type Create struct {
	Name     string          `json:"name" validate:"required"`
	UserName *string         `form:"username" validate:"required_without=Email"`
	Email    *string         `form:"email" validate:"required_without=UserName,omitempty,email"`
	Password string          `json:"password" validate:"required"`
	Role     enum.Permission `json:"role" validate:"required,oneof=admin student"`
}

func (c *Create) ToModel(hashedPassword string) *model.User {
	var userName *string
	if c.UserName != nil && *c.UserName != "" {
		userName = c.UserName
	}

	var email *string
	if c.Email != nil && *c.Email != "" {
		email = c.Email
	}

	return &model.User{
		Name:     c.Name,
		UserName: userName,
		Email:    email,
		Password: hashedPassword,
		Role:     c.Role,
	}
}

type Update struct {
	Name     string          `json:"name" validate:"required"`
	UserName *string         `json:"username" validate:"omitempty"`
	Email    *string         `json:"email" validate:"omitempty,email"`
	Password *string         `json:"password" validate:"omitempty,min=6"`
	Role     enum.Permission `json:"role" validate:"required,oneof=admin student"`
}
