package configure

import (
	"context"
	"fmt"

	"github.com/larek-tech/storage/redis"
)

func NewRedis(dsn string) *redis.Client {
	ctx := context.Background()

	cfg, err := redis.NewCfgFromDSN(dsn)
	if err != nil {
		panic(fmt.Errorf("parse redis dsn: %w", err))
	}

	client, err := redis.New(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("connect redis: %w", err))
	}
	return client
}
