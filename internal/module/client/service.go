package client

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, input *Create) error
	Update(ctx context.Context, id uuid.UUID, input *Update) error
	FindByID(ctx context.Context, id uuid.UUID) (*response.Client, error)
	FindAll(ctx context.Context) ([]response.Client, int64, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID, filter ClientFilter) (pagination.Response[response.Client], error)
	GetClient(ctx context.Context, data string) (*response.Client, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, input *Create) error {
	exists, err := s.repo.ExistCiByUser(ctx, userID, input.Ci)
	if err != nil {
		return err
	}
	if exists {
		return response.Conflict("Ya registraste un cliente con este C.I.")
	}

	client := input.ToModel(userID)

	if err := s.repo.Create(client); err != nil {
		return err
	}
	return nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input *Update) error {
	if err := s.repo.Exist(id); err != nil {
		return err
	}

	client, err := s.repo.FindByID(id)
	if err != nil {
		return response.NotFound("Cliente no encontrado")
	}

	input.ApplyTo(client)

	if err := s.repo.Update(client); err != nil {
		return err
	}
	return nil
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*response.Client, error) {
	client, err := s.repo.FindByID(id)
	if err != nil {
		return nil, response.NotFound("Cliente no encontrado")
	}
	return response.ClientToResponse(client), nil
}

func (s *service) FindAll(ctx context.Context) ([]response.Client, int64, error) {
	clients, total, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, 0, err
	}

	items := response.ClientsToResponse(clients)
	return items, total, nil
}

func (s *service) FindAllByUserID(ctx context.Context, userID uuid.UUID, filter ClientFilter) (pagination.Response[response.Client], error) {
	clients, total, err := s.repo.FindAllByUserID(ctx, userID, filter)
	if err != nil {
		return pagination.Response[response.Client]{}, err
	}

	items := response.ClientsToResponse(clients)
	return pagination.NewResponse(items, total, filter.Params), nil
}

func (s *service) GetClient(ctx context.Context, data string) (*response.Client, error) {
	client, err := s.repo.GetClient(data)
	if err != nil {
		return nil, response.NotFound("Cliente no encontrado")
	}
	return response.ClientToResponse(client), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Exist(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
