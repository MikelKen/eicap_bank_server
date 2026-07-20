package response

import (
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/model"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`

	Token string `json:"token,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func UsersToResponse(users []model.User) []User {
	items := make([]User, len(users))

	for i := range users {
		items[i] = *UserToResponse(&users[i], nil)
	}

	return items
}

func UserToResponse(user *model.User, token *string) *User {
	var deletedAt *time.Time
	if user.DeletedAt.Valid {
		deletedAt = &user.DeletedAt.Time
	}

	var tokenValue string
	if token != nil {
		tokenValue = *token
	}

	return &User{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     ifPtr(user.Email),
		Role:      user.Role.String(),
		Token:     tokenValue,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func ifPtr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
