package shop

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
	"github.com/EgorTarasov/shopx/internal/event"
)

type AuthClient interface {
	Register(ctx context.Context, in RegisterInput) (RegisterOutput, error)
	Login(ctx context.Context, in LoginInput) (LoginOutput, error)
	Validate(ctx context.Context, token domain.SessionToken) (MeOutput, error)
}

type OrderRepo interface {
	Create(ctx context.Context, o domain.Order) error
	GetByID(ctx context.Context, id domain.OrderID) (domain.Order, error)
}

type OrderEventPublisher interface {
	PublishCreated(ctx context.Context, e event.OrderCreatedEvent) error
}

type OrderIDGenerator interface {
	NewOrderID() domain.OrderID
}
