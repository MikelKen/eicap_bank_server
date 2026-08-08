package client

import (
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
)

type ClientFilter struct {
	pagination.Params
	Search string `form:"search" json:"search" query:"search"`
}

func (f *ClientFilter) SetDefaults() {
	f.Params.SetDefaults()
}

type Create struct {
	Name      string `json:"name" validate:"required"`
	Ci        string `json:"ci" validate:"required"`
	Sex       string `json:"sex" validate:"required,oneof=M F"`
	BirthDate string `json:"birth_date" validate:"required"`
}

func (c *Create) ToModel(userID uuid.UUID) *model.Client {
	return &model.Client{
		Name:      c.Name,
		Ci:        c.Ci,
		Sex:       c.Sex,
		BirthDate: c.BirthDate,
		UserID:    userID,
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
