package cashcount

type CountInput struct {
	DenominationID string `json:"denomination_id" validate:"required,uuid"`
	Quantity       int    `json:"quantity" validate:"required,gt=0"`
}
