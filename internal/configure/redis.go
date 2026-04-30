package configure

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/larek-tech/storage/redis"
)

func NewRedis(addr string) *redis.Client {
	ctx := context.Background()

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		panic(fmt.Errorf("parse redis addr %q: %w", addr, err))
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Errorf("parse redis port %q: %w", portStr, err))
	}

	cfg := redis.Cfg{Host: host, Port: port}
	client, err := redis.New(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("connect redis: %w", err))
	}
	return client
}
