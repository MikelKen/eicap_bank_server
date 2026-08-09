package typeaccount

import (
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
)

type TypeAccountFilter struct {
	pagination.Params
	Name string `form:"name" json:"name" query:"name"`
}

func (f *TypeAccountFilter) SetDefaults() {
	f.Params.SetDefaults()
}

type Create struct {
	Name string `json:"name" validate:"required"`
}

func (c *Create) ToModel() *model.TypeAccount {
	return &model.TypeAccount{
		Name: c.Name,
	}
}

type Update struct {
	Name string `json:"name" validate:"required"`
}

func (u *Update) ApplyTo(typeAccount *model.TypeAccount) {
	typeAccount.Name = u.Name
}
