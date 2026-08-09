package typeoperation

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input *Create) error
	Update(ctx context.Context, id uuid.UUID, input *Update) error
	FindByID(ctx context.Context, id uuid.UUID) (*response.TypeOperation, error)
	FindAll(ctx context.Context, filter TypeOperationFilter) (pagination.Response[response.TypeOperation], error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input *Create) error {
	typeOperation := input.ToModel()

	if err := s.repo.Create(typeOperation); err != nil {
		return err
	}
	return nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input *Update) error {
	if err := s.repo.Exist(id); err != nil {
		return err
	}

	typeOperation, err := s.repo.FindByID(id)
	if err != nil {
		return response.NotFound("Tipo de operación no encontrado")
	}

	input.ApplyTo(typeOperation)

	if err := s.repo.Update(typeOperation); err != nil {
		return err
	}
	return nil
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*response.TypeOperation, error) {
	typeOperation, err := s.repo.FindByID(id)
	if err != nil {
		return nil, response.NotFound("Tipo de operación no encontrado")
	}
	return response.TypeOperationToResponse(typeOperation), nil
}

func (s *service) FindAll(ctx context.Context, filter TypeOperationFilter) (pagination.Response[response.TypeOperation], error) {
	typeOperations, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return pagination.Response[response.TypeOperation]{}, err
	}

	items := response.TypeOperationsToResponse(typeOperations)
	return pagination.NewResponse(items, total, filter.Params), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Exist(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
