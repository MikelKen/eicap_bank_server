package user

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input *Create) error
	FindByID(ctx context.Context, userID uuid.UUID) (*response.User, error)
	FindAll(ctx context.Context, filter UserFilter) (pagination.Response[response.User], error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input *Create) error {
	if input.Email != nil && *input.Email != "" {
		if _, err := s.repo.GetByEmail(*input.Email); err == nil {
			return response.Conflict("El correo ya está en uso")
		}
	}

	hashedPassword, err := pkg.Hash(input.Password)
	if err != nil {
		return response.InternalServerError("Error al procesar la contraseña")
	}

	user := input.ToModel(hashedPassword)

	if err := s.repo.Create(user); err != nil {
		return err
	}
	return nil
}

func (s *service) FindByID(ctx context.Context, userID uuid.UUID) (*response.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return response.UserToResponse(user, nil), nil
}

func (s *service) FindAll(ctx context.Context, filter UserFilter) (pagination.Response[response.User], error) {
	users, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return pagination.Response[response.User]{}, err
	}

	items := response.UsersToResponse(users)
	return pagination.NewResponse(items, total, filter.Params), nil
}
