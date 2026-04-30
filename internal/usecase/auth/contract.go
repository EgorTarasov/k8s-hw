package auth

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type UserRepo interface {
	Create(ctx context.Context, u domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id domain.UserID) (domain.User, error)
}

type SessionRepo interface {
	Save(ctx context.Context, s domain.Session) error
	Get(ctx context.Context, token domain.SessionToken) (domain.Session, error)
	Delete(ctx context.Context, token domain.SessionToken) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenGenerator interface {
	New() domain.SessionToken
}

type UserIDGenerator interface {
	NewUserID() domain.UserID
}
