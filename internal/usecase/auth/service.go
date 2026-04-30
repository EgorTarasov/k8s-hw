package auth

import "time"

const minPasswordLen = 6

type Service struct {
	users    UserRepo
	sessions SessionRepo
	hasher   PasswordHasher
	tokens   TokenGenerator
	ids      UserIDGenerator
	now      func() time.Time
}

func NewService(users UserRepo, sessions SessionRepo, hasher PasswordHasher, tokens TokenGenerator, ids UserIDGenerator) *Service {
	return &Service{
		users:    users,
		sessions: sessions,
		hasher:   hasher,
		tokens:   tokens,
		ids:      ids,
		now:      time.Now,
	}
}
