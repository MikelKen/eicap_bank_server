package denomination

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input *Create) error
	Update(ctx context.Context, id uuid.UUID, input *Update) error
	FindByID(ctx context.Context, id uuid.UUID) (*response.Denomination, error)
	FindAll(ctx context.Context, filter DenominationFilter) (pagination.Response[response.Denomination], error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, input *Create) error {
	denomination := input.ToModel()

	if err := s.repo.Create(denomination); err != nil {
		return err
	}
	return nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input *Update) error {
	if err := s.repo.Exist(id); err != nil {
		return err
	}

	denomination, err := s.repo.FindByID(id)
	if err != nil {
		return response.NotFound("Denominación no encontrada")
	}

	input.ApplyTo(denomination)

	if err := s.repo.Update(denomination); err != nil {
		return err
	}
	return nil
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*response.Denomination, error) {
	denomination, err := s.repo.FindByID(id)
	if err != nil {
		return nil, response.NotFound("Denominación no encontrada")
	}
	return response.DenominationToResponse(denomination), nil
}

func (s *service) FindAll(ctx context.Context, filter DenominationFilter) (pagination.Response[response.Denomination], error) {
	denominations, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return pagination.Response[response.Denomination]{}, err
	}

	items := response.DenominationsToResponse(denominations)
	return pagination.NewResponse(items, total, filter.Params), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Exist(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
