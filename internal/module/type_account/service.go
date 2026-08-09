package typeaccount

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input *Create) error
	Update(ctx context.Context, id uuid.UUID, input *Update) error
	FindByID(ctx context.Context, id uuid.UUID) (*response.TypeAccount, error)
	FindAll(ctx context.Context, filter TypeAccountFilter) (pagination.Response[response.TypeAccount], error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input *Create) error {
	typeAccount := input.ToModel()

	if err := s.repo.Create(typeAccount); err != nil {
		return err
	}
	return nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input *Update) error {
	if err := s.repo.Exist(id); err != nil {
		return err
	}

	typeAccount, err := s.repo.FindByID(id)
	if err != nil {
		return response.NotFound("Tipo de cuenta no encontrado")
	}

	input.ApplyTo(typeAccount)

	if err := s.repo.Update(typeAccount); err != nil {
		return err
	}
	return nil
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*response.TypeAccount, error) {
	typeAccount, err := s.repo.FindByID(id)
	if err != nil {
		return nil, response.NotFound("Tipo de cuenta no encontrado")
	}
	return response.TypeAccountToResponse(typeAccount), nil
}

func (s *service) FindAll(ctx context.Context, filter TypeAccountFilter) (pagination.Response[response.TypeAccount], error) {
	typeAccounts, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return pagination.Response[response.TypeAccount]{}, err
	}

	items := response.TypeAccountsToResponse(typeAccounts)
	return pagination.NewResponse(items, total, filter.Params), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Exist(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
