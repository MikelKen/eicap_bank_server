package response

import (
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
)

type Client struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Ci        string `json:"ci"`
	Sex       string `json:"sex"`
	BirthDate string `json:"birth_date"`
	UserID    string `json:"user_id"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func ClientsToResponse(clients []model.Client) []Client {
	items := make([]Client, len(clients))

	for i := range clients {
		items[i] = *ClientToResponse(&clients[i])
	}

	return items
}

func ClientToResponse(client *model.Client) *Client {
	var deletedAt *time.Time
	if client.DeletedAt.Valid {
		deletedAt = &client.DeletedAt.Time
	}

	return &Client{
		ID:        client.ID.String(),
		Name:      client.Name,
		Ci:        client.Ci,
		Sex:       client.Sex,
		BirthDate: client.BirthDate,
		UserID:    client.UserID.String(),
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
