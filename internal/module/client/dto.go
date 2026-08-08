package client

import (
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
)

type ClientFilter struct {
	pagination.Params
	Name string `form:"name" json:"name" query:"name"`
	Ci   string `form:"ci" json:"ci" query:"ci"`
}

func (f *ClientFilter) SetDefaults() {
	f.Params.SetDefaults()
}

type Create struct {
	Name      string    `json:"name" validate:"required"`
	Ci        string    `json:"ci" validate:"required"`
	Sex       string    `json:"sex" validate:"required,oneof=M F"`
	BirthDate string    `json:"birth_date" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
}

func (c *Create) ToModel() *model.Client {
	return &model.Client{
		Name:      c.Name,
		Ci:        c.Ci,
		Sex:       c.Sex,
		BirthDate: c.BirthDate,
		UserID:    c.UserID,
	}
}

type Update struct {
	Name      string `json:"name" validate:"required"`
	Ci        string `json:"ci" validate:"required"`
	Sex       string `json:"sex" validate:"required,oneof=M F"`
	BirthDate string `json:"birth_date" validate:"required"`
}

func (u *Update) ApplyTo(client *model.Client) {
	client.Name = u.Name
	client.Ci = u.Ci
	client.Sex = u.Sex
	client.BirthDate = u.BirthDate
}
