package domain

import "time"

type OrderStatus string

const (
	StatusNew        OrderStatus = "new"
	StatusProcessing OrderStatus = "processing"
	StatusProcessed  OrderStatus = "processed"
	StatusFailed     OrderStatus = "failed"
)

type SKU string

func (s SKU) String() string { return string(s) }

type OrderItem struct {
	SKU SKU `json:"sku"`
	Qty int `json:"qty"`
}

type OrderID string

func (id OrderID) String() string { return string(id) }

type Order struct {
	ID        OrderID
	UserID    UserID
	Status    OrderStatus
	Items     []OrderItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s OrderStatus) CanTransitionTo(next OrderStatus) bool {
	switch s {
	case StatusNew:
		return next == StatusProcessing
	case StatusProcessing:
		return next == StatusProcessed || next == StatusFailed
	}
	return false
}
