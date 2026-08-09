package response

import (
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
)

type OperationInformation struct {
	ID          string `json:"id"`
	Origin      string `json:"origin"`
	Reason      string `json:"reason"`
	Destination string `json:"destination"`
	Details     string `json:"details"`
}

func OperationInformationToResponse(o *model.OperationInformation) *OperationInformation {
	if o == nil {
		return nil
	}
	return &OperationInformation{
		ID:          o.ID.String(),
		Origin:      o.Origin,
		Reason:      o.Reason,
		Destination: o.Destination,
		Details:     o.Details,
	}
}

type BankOperation struct {
	ID                string                `json:"id"`
	Code              string                `json:"code"`
	Date              time.Time             `json:"date"`
	PreviousBalance   string                `json:"previous_balance"`
	Import            string                `json:"import"`
	EndBalance        string                `json:"end_balance"`
	TypeOperationID   string                `json:"type_operation_id"`
	TypeOperationCode string                `json:"type_operation_code,omitempty"`
	AccountID         string                `json:"account_id,omitempty"`
	AccountNumber     string                `json:"account_number,omitempty"`
	CashSessionID     string                `json:"cash_session_id,omitempty"`
	Info              *OperationInformation `json:"info,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func BankOperationToResponse(b *model.BankOperation) *BankOperation {
	resp := &BankOperation{
		ID:                b.ID.String(),
		Code:              b.Code,
		Date:              b.Date,
		PreviousBalance:   b.PreviousBalance.StringFixed(2),
		Import:            b.Import.StringFixed(2),
		EndBalance:        b.EndBalance.StringFixed(2),
		TypeOperationID:   b.TypeOperationID.String(),
		TypeOperationCode: b.TypeOperation.Code,
		Info:              OperationInformationToResponse(b.OperationInformation),
		CreatedAt:         b.CreatedAt,
	}

	if b.AccountID != nil {
		resp.AccountID = b.AccountID.String()
	}
	if b.Account != nil {
		resp.AccountNumber = b.Account.Number
	}
	if b.CashSessionID != nil {
		resp.CashSessionID = b.CashSessionID.String()
	}

	return resp
}

func BankOperationsToResponse(list []model.BankOperation) []BankOperation {
	items := make([]BankOperation, len(list))
	for i := range list {
		items[i] = *BankOperationToResponse(&list[i])
	}
	return items
}
