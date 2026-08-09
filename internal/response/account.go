package response

import (
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
)

type Account struct {
	ID       string `json:"id"`
	Number   string `json:"number"`
	Interest string `json:"interest"`
	Balance  string `json:"balance"`
	Status   string `json:"status"`

	ClientID        string `json:"client_id"`
	ClientName      string `json:"client_name,omitempty"`
	TypeAccountID   string `json:"type_account_id"`
	TypeAccountName string `json:"type_account_name,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func AccountsToResponse(accounts []model.Account) []Account {
	items := make([]Account, len(accounts))

	for i := range accounts {
		items[i] = *AccountToResponse(&accounts[i])
	}

	return items
}

func AccountToResponse(account *model.Account) *Account {
	var deletedAt *time.Time
	if account.DeletedAt.Valid {
		deletedAt = &account.DeletedAt.Time
	}

	return &Account{
		ID:              account.ID.String(),
		Number:          account.Number,
		Interest:        account.Interest.StringFixed(2),
		Balance:         account.Balance.StringFixed(2),
		Status:          account.Status,
		ClientID:        account.ClientID.String(),
		ClientName:      account.Client.Name,
		TypeAccountID:   account.TypeAccountID.String(),
		TypeAccountName: account.TypeAccount.Name,
		CreatedAt:       account.CreatedAt,
		UpdatedAt:       account.UpdatedAt,
		DeletedAt:       deletedAt,
	}
}
