package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EgorTarasov/shopx/internal/domain"
	"github.com/larek-tech/storage/redis"
	goredis "github.com/redis/go-redis/v9"
)

type SessionRepo struct {
	client *redis.Client
}

func NewSessionRepo(client *redis.Client) *SessionRepo { return &SessionRepo{client: client} }

type sessionRecord struct {
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func sessionKey(token domain.SessionToken) string { return "session:" + token.String() }

func (r *SessionRepo) Save(ctx context.Context, s domain.Session) error {
	ttl := time.Until(s.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}
	rec := sessionRecord{UserID: s.UserID.String(), ExpiresAt: s.ExpiresAt}
	if err := r.client.SetJSON(ctx, sessionKey(s.Token), rec, ttl); err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}

func (r *SessionRepo) Get(ctx context.Context, token domain.SessionToken) (domain.Session, error) {
	var rec sessionRecord
	if err := r.client.GetJSON(ctx, sessionKey(token), &rec); err != nil {
		if errors.Is(err, goredis.Nil) {
			return domain.Session{}, domain.ErrSessionNotFound
		}
		return domain.Session{}, fmt.Errorf("load session: %w", err)
	}
	return domain.Session{
		Token:     token,
		UserID:    domain.UserID(rec.UserID),
		ExpiresAt: rec.ExpiresAt,
	}, nil
}

func (r *SessionRepo) Delete(ctx context.Context, token domain.SessionToken) error {
	if _, err := r.client.Del(ctx, sessionKey(token)); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}
