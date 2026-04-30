package shop

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type MeInput struct {
	Token domain.SessionToken
}

type MeOutput struct {
	ID    domain.UserID
	Email string
	Name  string
}

func (s *Service) Me(ctx context.Context, in MeInput) (MeOutput, error) {
	return s.auth.Validate(ctx, in.Token)
}
