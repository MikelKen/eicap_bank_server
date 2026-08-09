package denomination

import (
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/shopspring/decimal"
)

type DenominationFilter struct {
	pagination.Params
	Type string `form:"type" json:"type" query:"type"`
}

func (f *DenominationFilter) SetDefaults() {
	f.Params.SetDefaults()
}

type Create struct {
	Type  string          `json:"type" validate:"required,oneof=bill coin"`
	Value decimal.Decimal `json:"value" validate:"required"`
	Name  string          `json:"name" validate:"required"`
}

func (c *Create) ToModel() *model.Denomination {
	return &model.Denomination{
		Type:  c.Type,
		Value: c.Value,
		Name:  c.Name,
	}
}

type Update struct {
	Type  string          `json:"type" validate:"required,oneof=bill coin"`
	Value decimal.Decimal `json:"value" validate:"required"`
	Name  string          `json:"name" validate:"required"`
}

func (u *Update) ApplyTo(denomination *model.Denomination) {
	denomination.Type = u.Type
	denomination.Value = u.Value
	denomination.Name = u.Name
}
