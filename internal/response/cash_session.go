package response

import (
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
)

type CashSession struct {
	ID               string      `json:"id"`
	State            string      `json:"state"`
	OpeningDate      time.Time   `json:"opening_date"`
	OpeningAmount    string      `json:"opening_amount"`
	ClosingDate      *time.Time  `json:"closing_date,omitempty"`
	ClosingAmount    *string     `json:"closing_amount,omitempty"`
	ExpectedAmount   *string     `json:"expected_amount,omitempty"`
	DifferenceAmount *string     `json:"difference_amount,omitempty"`
	OperationCode    string      `json:"operation_code,omitempty"`
	UserID           string      `json:"user_id"`
	UserName         string      `json:"user_name,omitempty"`
	Counts           []CashCount `json:"counts,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func CashSessionsToResponse(sessions []model.CashSession) []CashSession {
	items := make([]CashSession, len(sessions))
	for i := range sessions {
		items[i] = *CashSessionToResponse(&sessions[i])
	}
	return items
}

func CashSessionToResponse(s *model.CashSession) *CashSession {
	var deletedAt *time.Time
	if s.DeletedAt.Valid {
		deletedAt = &s.DeletedAt.Time
	}

	resp := &CashSession{
		ID:            s.ID.String(),
		State:         s.State,
		OpeningDate:   s.OpeningDate,
		OpeningAmount: s.OpeningAmount.StringFixed(2),
		UserID:        s.UserID.String(),
		UserName:      s.User.Name,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
		DeletedAt:     deletedAt,
	}

	// Solo muestra datos de cierre si la sesión ya fue cerrada.
	if s.State == "closed" {
		closingDate := s.ClosingDate
		closingAmount := s.ClosingAmount.StringFixed(2)
		expectedAmount := s.ExpectedAmount.StringFixed(2)
		differenceAmount := s.DifferenceAmount.StringFixed(2)

		resp.ClosingDate = &closingDate
		resp.ClosingAmount = &closingAmount
		resp.ExpectedAmount = &expectedAmount
		resp.DifferenceAmount = &differenceAmount
	}

	if len(s.CashCounts) > 0 {
		resp.Counts = CashCountsToResponse(s.CashCounts)
	}

	// Código de la operación bancaria que dejó constancia del cierre de caja (CICA).
	for _, op := range s.BankOperations {
		if op.TypeOperation.Code == "CICA" {
			resp.OperationCode = op.Code
			break
		}
	}

	return resp
}
