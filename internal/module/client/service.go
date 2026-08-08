package client

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input *Create) error
	Update(ctx context.Context, id uuid.UUID, input *Update) error
	FindByID(ctx context.Context, id uuid.UUID) (*response.Client, error)
	FindAll(ctx context.Context) ([]response.Client, int64, error)
	GetClient(ctx context.Context, data string) (*response.Client, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input *Create) error {
	if _, err := s.repo.GetClient(input.Ci); err == nil {
		return response.Conflict("El C.I. ya está registrado")
	}

	client := input.ToModel()

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
