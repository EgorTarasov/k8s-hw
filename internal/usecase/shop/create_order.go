package shop

import (
	"context"
	"fmt"

	"github.com/EgorTarasov/shopx/internal/domain"
	"github.com/EgorTarasov/shopx/internal/event"
)

type OrderItemInput struct {
	SKU domain.SKU
	Qty int
}

type CreateOrderInput struct {
	UserID domain.UserID
	Items  []OrderItemInput
}

type CreateOrderOutput struct {
	ID     domain.OrderID
	Status domain.OrderStatus
}

func (s *Service) CreateOrder(ctx context.Context, in CreateOrderInput) (CreateOrderOutput, error) {
	if len(in.Items) == 0 {
		return CreateOrderOutput{}, fmt.Errorf("%w: items must not be empty", domain.ErrInvalidInput)
	}
	items := make([]domain.OrderItem, len(in.Items))
	for i, it := range in.Items {
		if it.SKU == "" || it.Qty <= 0 {
			return CreateOrderOutput{}, fmt.Errorf("%w: invalid item at index %d", domain.ErrInvalidInput, i)
		}
		items[i] = domain.OrderItem{SKU: it.SKU, Qty: it.Qty}
	}

	now := s.now().UTC()
	o := domain.Order{
		ID:        s.ids.NewOrderID(),
		UserID:    in.UserID,
		Status:    domain.StatusNew,
		Items:     items,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.orders.Create(ctx, o); err != nil {
		return CreateOrderOutput{}, fmt.Errorf("create order: %w", err)
	}

	evt := event.OrderCreatedEvent{
		ID:         o.ID,
		UserID:     o.UserID,
		Items:      o.Items,
		OccurredAt: now,
	}
	if err := s.events.PublishCreated(ctx, evt); err != nil {
		s.logger.Error("publish order.created failed", "order_id", o.ID, "err", err)
	}

	return CreateOrderOutput{ID: o.ID, Status: o.Status}, nil
}
