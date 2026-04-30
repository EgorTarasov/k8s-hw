package auth

import (
	"context"
	"fmt"
	"strings"

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
	email := strings.ToLower(strings.TrimSpace(in.Email))
	name := strings.TrimSpace(in.Name)
	if email == "" || name == "" {
		return RegisterOutput{}, fmt.Errorf("%w: email and name are required", domain.ErrInvalidInput)
	}
	if len(in.Password) < minPasswordLen {
		return RegisterOutput{}, fmt.Errorf("%w: password must be at least %d characters", domain.ErrInvalidInput, minPasswordLen)
	}

	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return RegisterOutput{}, fmt.Errorf("hash password: %w", err)
	}

	u := domain.User{
		ID:           s.ids.NewUserID(),
		Email:        email,
		Name:         name,
		PasswordHash: hash,
		CreatedAt:    s.now().UTC(),
	}
	if err := s.users.Create(ctx, u); err != nil {
		return RegisterOutput{}, err
	}
	return RegisterOutput{ID: u.ID, Email: u.Email, Name: u.Name}, nil
}
