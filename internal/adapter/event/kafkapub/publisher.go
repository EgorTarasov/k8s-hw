package kafkapub

import (
	"context"
	"fmt"

	"github.com/EgorTarasov/shopx/internal/event"
	"github.com/larek-tech/storage/kafka"
)

type Publisher struct {
	producer *kafka.Producer[event.OrderCreatedEvent]
}

func New(tr kafka.Transport, topic string) *Publisher {
	p := kafka.NewProducer[event.OrderCreatedEvent](tr, topic)
	return &Publisher{producer: p}
}

func (p *Publisher) PublishCreated(ctx context.Context, e event.OrderCreatedEvent) error {
	if err := p.producer.Produce(ctx, []event.OrderCreatedEvent{e}); err != nil {
		return fmt.Errorf("kafka produce: %w", err)
	}
	return nil
}
