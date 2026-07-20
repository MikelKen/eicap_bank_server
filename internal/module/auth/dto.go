package auth

type Login struct {
	Email    *string `json:"email" validate:"required_without=UserName,omitempty,email"`
	Password string  `json:"password" validate:"required"`
}
