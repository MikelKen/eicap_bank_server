package response

import (
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
)

type TypeAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func TypeAccountsToResponse(typeAccounts []model.TypeAccount) []TypeAccount {
	items := make([]TypeAccount, len(typeAccounts))

	for i := range typeAccounts {
		items[i] = *TypeAccountToResponse(&typeAccounts[i])
	}

	return items
}

func TypeAccountToResponse(typeAccount *model.TypeAccount) *TypeAccount {
	var deletedAt *time.Time
	if typeAccount.DeletedAt.Valid {
		deletedAt = &typeAccount.DeletedAt.Time
	}

	return &TypeAccount{
		ID:        typeAccount.ID.String(),
		Name:      typeAccount.Name,
		CreatedAt: typeAccount.CreatedAt,
		UpdatedAt: typeAccount.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
