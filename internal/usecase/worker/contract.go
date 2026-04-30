package worker

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type OrderRepo interface {
	GetByID(ctx context.Context, id domain.OrderID) (domain.Order, error)
	UpdateStatus(ctx context.Context, id domain.OrderID, status domain.OrderStatus) error
}

type OrderProcessor interface {
	Process(ctx context.Context, o domain.Order) error
}
