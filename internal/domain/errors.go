package domain

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailTaken      = errors.New("email already taken")
	ErrBadCredentials  = errors.New("bad credentials")
	ErrSessionNotFound = errors.New("session not found")
	ErrOrderNotFound   = errors.New("order not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidStatus   = errors.New("invalid status transition")
	ErrForbidden       = errors.New("forbidden")
)
