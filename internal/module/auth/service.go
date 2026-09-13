package auth

import (
	"errors"
	"time"

	"github.com/Eicap/EICAP-BANK/server/internal/config"
	"github.com/Eicap/EICAP-BANK/server/internal/model"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg"
	"gorm.io/gorm"
)

type Service interface {
	Login(input *Login) (string, time.Duration, *response.User, error)
}

type userRepo interface {
	GetByEmail(email string) (*model.User, error)
	GetByUserName(username string) (*model.User, error)
}

type service struct {
	userRepo userRepo
	cfg      *config.Config
}

func NewService(userRepo userRepo, cfg *config.Config) Service {
	return &service{userRepo: userRepo, cfg: cfg}
}

func (s *service) Login(input *Login) (string, time.Duration, *response.User, error) {
	var user *model.User
	var err error

	if input.Email != nil {
		user, err = s.userRepo.GetByEmail(*input.Email)
	} else if input.UserName != nil {
		user, err = s.userRepo.GetByUserName(*input.UserName)
	} else {
		return "", 0, nil, response.BadRequest("Email o username son requeridos")
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, nil, response.Unauthorized("Credenciales inválidas")
		}
		return "", 0, nil, response.InternalServerError("Error al buscar usuario")
	}

	if user == nil {
		return "", 0, nil, response.Unauthorized("Credenciales inválidas")
	}

	if err := pkg.Compare(input.Password, user.Password); err != nil {
		return "", 0, nil, response.Unauthorized("Credenciales inválidas")
	}

	expiration, err := time.ParseDuration(s.cfg.JWTExpiration)
	if err != nil {
		return "", 0, nil, response.InternalServerError("Error de configuración del servidor")
	}

	token, err := pkg.GenerateToken(user, s.cfg.JWTSecret, expiration)
	if err != nil {
		return "", 0, nil, response.InternalServerError("Error al generar token de acceso")
	}

	userResponse := response.UserToResponse(user, &token)

	return token, expiration, userResponse, nil
}
