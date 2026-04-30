package domain

import "time"

type UserID string

func (id UserID) String() string { return string(id) }

type User struct {
	ID           UserID
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    time.Time
}
