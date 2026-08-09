package account

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/module/client"
	typeaccount "github.com/Eicap/EICAP-BANK/server/internal/module/type_account"
	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/Eicap/EICAP-BANK/server/pkg/pagination"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input *Create) error
	Update(ctx context.Context, id uuid.UUID, input *Update) error
	FindByID(ctx context.Context, id uuid.UUID) (*response.Account, error)
	FindAll(ctx context.Context, filter AccountFilter) (pagination.Response[response.Account], error)
	FindAllByClientID(ctx context.Context, clientID uuid.UUID, filter AccountFilter) (pagination.Response[response.Account], error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo            Repo
	clientRepo      client.Repo
	typeAccountRepo typeaccount.Repo
}

func NewService(repo Repo, clientRepo client.Repo, typeAccountRepo typeaccount.Repo) Service {
	return &service{repo: repo, clientRepo: clientRepo, typeAccountRepo: typeAccountRepo}
}

func (s *service) Create(ctx context.Context, input *Create) error {
	if err := s.clientRepo.Exist(input.ClientID); err != nil {
		return response.NotFound("El cliente no existe")
	}

	if err := s.typeAccountRepo.Exist(input.TypeAccountID); err != nil {
		return response.NotFound("El tipo de cuenta no existe")
	}

	number, err := s.repo.NextNumber()
	if err != nil {
		return err
	}

	acc := input.ToModel(number)

	if err := s.repo.Create(acc); err != nil {
		return err
	}
	return nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input *Update) error {
	if err := s.repo.Exist(id); err != nil {
		return err
	}

	acc, err := s.repo.FindByID(id)
	if err != nil {
		return response.NotFound("Cuenta no encontrada")
	}

	input.ApplyTo(acc)

	if err := s.repo.Update(acc); err != nil {
		return err
	}
	return nil
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID) (*response.Account, error) {
	acc, err := s.repo.FindByID(id)
	if err != nil {
		return nil, response.NotFound("Cuenta no encontrada")
	}
	return response.AccountToResponse(acc), nil
}

func (s *service) FindAll(ctx context.Context, filter AccountFilter) (pagination.Response[response.Account], error) {
	accounts, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return pagination.Response[response.Account]{}, err
	}

	items := response.AccountsToResponse(accounts)
	return pagination.NewResponse(items, total, filter.Params), nil
}

func (s *service) FindAllByClientID(ctx context.Context, clientID uuid.UUID, filter AccountFilter) (pagination.Response[response.Account], error) {
	accounts, total, err := s.repo.FindAllByClientID(ctx, clientID, filter)
	if err != nil {
		return pagination.Response[response.Account]{}, err
	}

	items := response.AccountsToResponse(accounts)
	return pagination.NewResponse(items, total, filter.Params), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Exist(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
