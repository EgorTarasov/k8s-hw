package configure

import (
	"context"
	"fmt"
	"strings"

	"github.com/larek-tech/storage/kafka"
)

func NewKafkaTransport(brokers, groupID string) *kafka.KafkaGoTransport {
	ctx := context.Background()

	cfg := kafka.Cfg{
		Brokers:  splitBrokers(brokers),
		GroupID:  groupID,
		ClientID: "shopx",
	}

	tr, err := kafka.NewKafkaGoTransport(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("init kafka transport: %w", err))
	}
	return tr
}

func splitBrokers(s string) []string {
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
