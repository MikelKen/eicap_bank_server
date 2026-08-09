package account

import (
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountFilter struct {
	pagination.Params
	Status        string `form:"status" json:"status" query:"status"`
	ClientID      string `form:"client_id" json:"client_id" query:"client_id"`
	TypeAccountID string `form:"type_account_id" json:"type_account_id" query:"type_account_id"`
}

func (f *AccountFilter) SetDefaults() {
	f.Params.SetDefaults()
}

type Create struct {
	ClientID      uuid.UUID       `json:"client_id" validate:"required"`
	TypeAccountID uuid.UUID       `json:"type_account_id" validate:"required"`
	Interest      decimal.Decimal `json:"interest"`
}

func (c *Create) ToModel(number string) *model.Account {
	return &model.Account{
		Number:        number,
		Interest:      c.Interest,
		Balance:       decimal.Zero,
		Status:        "active",
		ClientID:      c.ClientID,
		TypeAccountID: c.TypeAccountID,
	}
}

type Update struct {
	Status   string          `json:"status" validate:"required,oneof=active inactive blocked"`
	Interest decimal.Decimal `json:"interest"`
}

func (u *Update) ApplyTo(account *model.Account) {
	account.Status = u.Status
	account.Interest = u.Interest
}
