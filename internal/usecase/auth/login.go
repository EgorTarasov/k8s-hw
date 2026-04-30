package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token     domain.SessionToken
	ExpiresAt time.Time
}

func (s *Service) Login(ctx context.Context, in LoginInput) (LoginOutput, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	if email == "" || in.Password == "" {
		return LoginOutput{}, fmt.Errorf("%w: email and password are required", domain.ErrInvalidInput)
	}

	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return LoginOutput{}, domain.ErrBadCredentials
		}
		return LoginOutput{}, err
	}
	if err := s.hasher.Compare(u.PasswordHash, in.Password); err != nil {
		return LoginOutput{}, domain.ErrBadCredentials
	}

	token := s.tokens.New()
	expiresAt := s.now().Add(domain.SessionTTL).UTC()
	if err := s.sessions.Save(ctx, domain.Session{
		Token:     token,
		UserID:    u.ID,
		ExpiresAt: expiresAt,
	}); err != nil {
		return LoginOutput{}, fmt.Errorf("save session: %w", err)
	}
	return LoginOutput{Token: token, ExpiresAt: expiresAt}, nil
}
