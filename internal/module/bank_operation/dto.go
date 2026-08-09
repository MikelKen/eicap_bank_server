package bankoperation

import (
	"github.com/google/uuid"

	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/shopspring/decimal"
)

const (
	CodeIncome  = "ING" // Ingreso
	CodeExpense = "EGR" // Egreso
)

type BankOperationFilter struct {
	pagination.Params
	AccountID         string `form:"account_id" json:"account_id" query:"account_id"`
	TypeOperationCode string `form:"type_operation_code" json:"type_operation_code" query:"type_operation_code"`
}

func (f *BankOperationFilter) SetDefaults() {
	f.Params.SetDefaults()
}

// OperationInformationInput solo se exige cuando TypeOperationCode es ING o EGR.
type OperationInformationInput struct {
	Origin      string `json:"origin" validate:"required"`
	Reason      string `json:"reason" validate:"required"`
	Destination string `json:"destination" validate:"required"`
	Details     string `json:"details"`
}

type Create struct {
	TypeOperationCode string                     `json:"type_operation_code" validate:"required,oneof=ING EGR"`
	AccountID         uuid.UUID                  `json:"account_id" validate:"required"`
	Amount            decimal.Decimal            `json:"amount" validate:"required"`
	Info              *OperationInformationInput `json:"info" validate:"required"`
}
