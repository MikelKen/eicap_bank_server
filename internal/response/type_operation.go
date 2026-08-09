package response

import (
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
)

type TypeOperation struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func TypeOperationsToResponse(typeOperations []model.TypeOperation) []TypeOperation {
	items := make([]TypeOperation, len(typeOperations))

	for i := range typeOperations {
		items[i] = *TypeOperationToResponse(&typeOperations[i])
	}

	return items
}

func TypeOperationToResponse(typeOperation *model.TypeOperation) *TypeOperation {
	var deletedAt *time.Time
	if typeOperation.DeletedAt.Valid {
		deletedAt = &typeOperation.DeletedAt.Time
	}

	return &TypeOperation{
		ID:        typeOperation.ID.String(),
		Code:      typeOperation.Code,
		Name:      typeOperation.Name,
		CreatedAt: typeOperation.CreatedAt,
		UpdatedAt: typeOperation.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
