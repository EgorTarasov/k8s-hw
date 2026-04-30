package auth

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type ValidateInput struct {
	Token domain.SessionToken
}

type ValidateOutput struct {
	ID    domain.UserID
	Email string
	Name  string
}

func (s *Service) Validate(ctx context.Context, in ValidateInput) (ValidateOutput, error) {
	if in.Token == "" {
		return ValidateOutput{}, domain.ErrSessionNotFound
	}
	sess, err := s.sessions.Get(ctx, in.Token)
	if err != nil {
		return ValidateOutput{}, err
	}
	if !s.now().Before(sess.ExpiresAt) {
		_ = s.sessions.Delete(ctx, in.Token)
		return ValidateOutput{}, domain.ErrSessionNotFound
	}
	u, err := s.users.GetByID(ctx, sess.UserID)
	if err != nil {
		return ValidateOutput{}, err
	}
	return ValidateOutput{ID: u.ID, Email: u.Email, Name: u.Name}, nil
}
