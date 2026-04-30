package shop

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token domain.SessionToken
}

func (s *Service) Login(ctx context.Context, in LoginInput) (LoginOutput, error) {
	return s.auth.Login(ctx, in)
}
