package response

import "github.com/Eicap/EICAP-BANK/server/internal/model"

type CashCount struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	Quantity          int    `json:"quantity"`
	Subtotal          string `json:"subtotal"`
	DenominationID    string `json:"denomination_id"`
	DenominationName  string `json:"denomination_name,omitempty"`
	DenominationValue string `json:"denomination_value,omitempty"`
}

func CashCountsToResponse(counts []model.CashCount) []CashCount {
	items := make([]CashCount, len(counts))
	for i := range counts {
		items[i] = *CashCountToResponse(&counts[i])
	}
	return items
}

func CashCountToResponse(c *model.CashCount) *CashCount {
	return &CashCount{
		ID:                c.ID.String(),
		Type:              c.Type,
		Quantity:          c.Quantity,
		Subtotal:          c.Subtotal.StringFixed(2),
		DenominationID:    c.DenominationID.String(),
		DenominationName:  c.Denomination.Name,
		DenominationValue: c.Denomination.Value.StringFixed(2),
	}
}
