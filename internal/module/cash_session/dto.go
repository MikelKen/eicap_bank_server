package cashsession

import (
	cashcount "github.com/Eicap/EICAP-BANK/server/internal/module/cash_count"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
)

const (
	StateOpen   = "open"
	StateClosed = "closed"
)

const (
	CountTypeOpening = "opening"
	CountTypeClosing = "closing"
)

type CashSessionFilter struct {
	pagination.Params
	State string `form:"state" json:"state" query:"state"`
}

func (f *CashSessionFilter) SetDefaults() {
	f.Params.SetDefaults()
}

// Open es lo que el cajero manda al abrir caja: el detalle de billetes/monedas contados.
type Open struct {
	Counts []cashcount.CountInput `json:"counts" validate:"required,min=1,dive"`
}

// Close es lo que el cajero manda al cerrar caja.
type Close struct {
	Counts []cashcount.CountInput `json:"counts" validate:"required,min=1,dive"`
}
