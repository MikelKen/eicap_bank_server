package typeoperation

import (
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
)

type TypeOperationFilter struct {
	pagination.Params
	Name string `form:"name" json:"name" query:"name"`
	Code string `form:"code" json:"code" query:"code"`
}

func (f *TypeOperationFilter) SetDefaults() {
	f.Params.SetDefaults()
}

type Create struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
}

func (c *Create) ToModel() *model.TypeOperation {
	return &model.TypeOperation{
		Name: c.Name,
		Code: c.Code,
	}
}

type Update struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
}

func (u *Update) ApplyTo(typeOperation *model.TypeOperation) {
	typeOperation.Name = u.Name
	typeOperation.Code = u.Code
}
