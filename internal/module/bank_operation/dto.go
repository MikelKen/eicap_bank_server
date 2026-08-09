package bankoperation

import (
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	CodeIncome      = "ING"  // Ingreso
	CodeExpense     = "EGR"  // Egreso
	CodeAccountOpen = "APC"  // Apertura de Cuenta
	CodeCashOpen    = "APCA" // Apertura de Caja
	CodeCashClose   = "CICA" // Cierre de Caja
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
// Para el resto de operaciones (APC, APCA, CICA...) no se registra información.
type OperationInformationInput struct {
	Origin      string `json:"origin"`
	Reason      string `json:"reason"`
	Destination string `json:"destination"`
	Details     string `json:"details"`
}

type Create struct {
	TypeOperationCode string                     `json:"type_operation_code" validate:"required,oneof=ING EGR APC APCA CICA"`
	AccountID         *uuid.UUID                 `json:"account_id"`
	Amount            decimal.Decimal            `json:"amount"`
	Info              *OperationInformationInput `json:"info"`
}
