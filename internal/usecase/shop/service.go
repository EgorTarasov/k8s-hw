package shop

import (
	"log/slog"
	"time"
)

type Service struct {
	auth   AuthClient
	orders OrderRepo
	events OrderEventPublisher
	ids    OrderIDGenerator
	logger *slog.Logger
	now    func() time.Time
}

func NewService(auth AuthClient, orders OrderRepo, events OrderEventPublisher, ids OrderIDGenerator, logger *slog.Logger) *Service {
	return &Service{
		auth:   auth,
		orders: orders,
		events: events,
		ids:    ids,
		logger: logger,
		now:    time.Now,
	}
}
