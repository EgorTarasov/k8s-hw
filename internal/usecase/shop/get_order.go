package shop

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type GetOrderInput struct {
	UserID domain.UserID
	ID     domain.OrderID
}

type GetOrderOutput struct {
	ID     domain.OrderID
	UserID domain.UserID
	Status domain.OrderStatus
	Items  []OrderItemInput
}

func (s *Service) GetOrder(ctx context.Context, in GetOrderInput) (GetOrderOutput, error) {
	o, err := s.orders.GetByID(ctx, in.ID)
	if err != nil {
		return GetOrderOutput{}, err
	}
	if o.UserID != in.UserID {
		return GetOrderOutput{}, domain.ErrOrderNotFound
	}
	items := make([]OrderItemInput, len(o.Items))
	for i, it := range o.Items {
		items[i] = OrderItemInput{SKU: it.SKU, Qty: it.Qty}
	}
	return GetOrderOutput{
		ID:     o.ID,
		UserID: o.UserID,
		Status: o.Status,
		Items:  items,
	}, nil
}
