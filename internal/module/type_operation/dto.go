package typeoperation

import (
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
)

const (
	CashFlowIncome  = "ING"
	CashFlowExpense = "EGR"
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
	Name         string `json:"name" validate:"required"`
	Code         string `json:"code" validate:"required"`
	CashFlowType string `json:"cash_flow_type" validate:"omitempty,oneof=ING EGR"`
}

func (c *Create) ToModel() *model.TypeOperation {
	return &model.TypeOperation{
		Name:         c.Name,
		Code:         c.Code,
		CashFlowType: c.CashFlowType,
	}
}

type Update struct {
	Name         string `json:"name" validate:"required"`
	Code         string `json:"code" validate:"required"`
	CashFlowType string `json:"cash_flow_type" validate:"omitempty,oneof=ING EGR"`
}

func (u *Update) ApplyTo(typeOperation *model.TypeOperation) {
	typeOperation.Name = u.Name
	typeOperation.Code = u.Code
	typeOperation.CashFlowType = u.CashFlowType
}
