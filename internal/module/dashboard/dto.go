package dashboard

import "time"

// NameCount es un par etiqueta/valor para distribuciones simples (roles...).
type NameCount struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

// TypeCount es el conteo y volumen por tipo de operación bancaria.
type TypeCount struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
	Total string `json:"total"`
}

type ClientsStats struct {
	Total        int64 `json:"total"`
	NewThisMonth int64 `json:"new_this_month"`
}

type OperationsStats struct {
	Total        int64       `json:"total"`
	IncomeTotal  string      `json:"income_total"`
	ExpenseTotal string      `json:"expense_total"`
	NetTotal     string      `json:"net_total"`
	ByType       []TypeCount `json:"by_type"`
}

type CashSessionsStats struct {
	Openings      int64 `json:"openings"`
	Closings      int64 `json:"closings"`
	CurrentlyOpen int64 `json:"currently_open"`
}

// ActivityPoint es un punto de la serie diaria de actividad (operaciones).
type ActivityPoint struct {
	Date       string `json:"date"` // YYYY-MM-DD
	Operations int64  `json:"operations"`
	Income     string `json:"income"`
	Expense    string `json:"expense"`
}

type RecentOperation struct {
	ID       string    `json:"id"`
	Code     string    `json:"code"`
	Date     time.Time `json:"date"`
	TypeCode string    `json:"type_code"`
	TypeName string    `json:"type_name"`
	Import   string    `json:"import"`
}

type UserPerformance struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Clients    int64     `json:"clients"`
	Sessions   int64     `json:"sessions"`
	Operations int64     `json:"operations"`
	CreatedAt  time.Time `json:"created_at"`
}

// AdminStats solo se incluye para usuarios admin.
type AdminStats struct {
	UsersTotal            int64             `json:"users_total"`
	UsersCreatedThisMonth int64             `json:"users_created_this_month"`
	UsersByRole           []NameCount       `json:"users_by_role"`
	TopUsers              []UserPerformance `json:"top_users"`
}

// Summary es la respuesta completa del panel.
type Summary struct {
	Role             string            `json:"role"`
	Clients          ClientsStats      `json:"clients"`
	Operations       OperationsStats   `json:"operations"`
	CashSessions     CashSessionsStats `json:"cash_sessions"`
	Activity         []ActivityPoint   `json:"activity"`
	RecentOperations []RecentOperation `json:"recent_operations"`
	Admin            *AdminStats       `json:"admin,omitempty"`
}
