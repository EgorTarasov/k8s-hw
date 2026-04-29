package configure

import (
	"context"
	"fmt"

	"github.com/larek-tech/storage/kafka"
)

func NewKafkaTransport(dsn string) *kafka.KafkaGoTransport {
	ctx := context.Background()

	cfg, err := kafka.NewCfgFromDSN(dsn)
	if err != nil {
		panic(fmt.Errorf("parse kafka dsn: %w", err))
	}

	tr, err := kafka.NewKafkaGoTransport(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("init kafka transport: %w", err))
	}
	return tr
}
