package cashcount

import (
	"context"

	"github.com/Eicap/EICAP-BANK/server/internal/response"
	"github.com/google/uuid"
)

type Service interface {
	FindAllBySessionID(ctx context.Context, sessionID uuid.UUID) ([]response.CashCount, error)
}

type service struct {
	repo Repo
}

func NewService(repo Repo) Service {
	return &service{repo: repo}
}

func (s *service) FindAllBySessionID(ctx context.Context, sessionID uuid.UUID) ([]response.CashCount, error) {
	counts, err := s.repo.FindAllBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return response.CashCountsToResponse(counts), nil
}
