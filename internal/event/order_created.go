package event

import (
	"time"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type OrderCreatedEvent struct {
	ID         domain.OrderID     `json:"id"`
	UserID     domain.UserID      `json:"userId"`
	Items      []domain.OrderItem `json:"items"`
	OccurredAt time.Time          `json:"occurredAt"`
}
