package domain

import "time"

const SessionTTL = 24 * time.Hour

type SessionToken string

func (t SessionToken) String() string { return string(t) }

type Session struct {
	Token     SessionToken
	UserID    UserID
	ExpiresAt time.Time
}
