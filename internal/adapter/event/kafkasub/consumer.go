package kafkasub

import (
	"context"
	"fmt"

	"github.com/EgorTarasov/shopx/internal/event"
	"github.com/larek-tech/storage/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer[event.OrderCreatedEvent]
}

func New(tr kafka.Transport, topic string) *Consumer {
	c := kafka.NewConsumer[event.OrderCreatedEvent](tr, topic)
	return &Consumer{consumer: c}
}

type Handler func(ctx context.Context, e event.OrderCreatedEvent) error

func (c *Consumer) Run(ctx context.Context, handler Handler) error {
	return c.consumer.Consume(ctx, func(ctx context.Context, msgs []kafka.Message[event.OrderCreatedEvent]) (int, error) {
		for i, m := range msgs {
			if err := handler(ctx, m.Value); err != nil {
				return i, fmt.Errorf("handle order.created: %w", err)
			}
		}
		return len(msgs), nil
	})
}
