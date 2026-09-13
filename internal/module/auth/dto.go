package auth

type Login struct {
	Email    *string `json:"email" validate:"required_without=UserName,omitempty,email"`
	UserName *string `json:"username" validate:"required_without=Email"`
	Password string  `json:"password" validate:"required"`
}
