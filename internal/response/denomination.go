package response

import (
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
)

type Denomination struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Name  string `json:"name"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func DenominationsToResponse(denominations []model.Denomination) []Denomination {
	items := make([]Denomination, len(denominations))

	for i := range denominations {
		items[i] = *DenominationToResponse(&denominations[i])
	}

	return items
}

func DenominationToResponse(denomination *model.Denomination) *Denomination {
	var deletedAt *time.Time
	if denomination.DeletedAt.Valid {
		deletedAt = &denomination.DeletedAt.Time
	}

	return &Denomination{
		ID:        denomination.ID.String(),
		Type:      denomination.Type,
		Value:     denomination.Value.StringFixed(2),
		Name:      denomination.Name,
		CreatedAt: denomination.CreatedAt,
		UpdatedAt: denomination.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
