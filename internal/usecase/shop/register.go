package shop

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type RegisterInput struct {
	Email    string
	Password string
	Name     string
}

type RegisterOutput struct {
	ID    domain.UserID
	Email string
	Name  string
}

func (s *Service) Register(ctx context.Context, in RegisterInput) (RegisterOutput, error) {
	return s.auth.Register(ctx, in)
}
