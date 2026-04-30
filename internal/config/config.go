package config

import (
	"flag"
	"fmt"
)

type AuthConfig struct {
	Listen      string
	PostgresDSN string
	RedisDSN    string
}

type ShopConfig struct {
	Listen      string
	AuthGRPC    string
	PostgresDSN string
	KafkaAddr   string
	OrdersTopic string
}

type WorkerConfig struct {
	PostgresDSN string
	KafkaAddr   string
	OrdersTopic string
	GroupID     string
	Workers     int
}

var (
	Auth   AuthConfig
	Shop   ShopConfig
	Worker WorkerConfig
)

func LoadAuth() error {
	flag.StringVar(&Auth.Listen, "listen", "127.0.0.1:9090", "gRPC listen address")
	flag.StringVar(&Auth.PostgresDSN, "postgres", "", "postgres DSN")
	flag.StringVar(&Auth.RedisDSN, "redis", "", "redis address (host:port)")
	flag.Parse()

	if Auth.PostgresDSN == "" {
		return fmt.Errorf("--postgres is required")
	}
	if Auth.RedisDSN == "" {
		return fmt.Errorf("--redis is required")
	}
	return nil
}

func LoadShop() error {
	flag.StringVar(&Shop.Listen, "listen", "127.0.0.1:8080", "HTTP listen address")
	flag.StringVar(&Shop.AuthGRPC, "auth-grpc", "", "auth-service gRPC address")
	flag.StringVar(&Shop.PostgresDSN, "postgres", "", "postgres DSN")
	flag.StringVar(&Shop.KafkaAddr, "kafka", "", "kafka brokers (comma-separated host:port)")
	flag.StringVar(&Shop.OrdersTopic, "orders-topic", "orders.created", "kafka topic for created orders")
	flag.Parse()

	if Shop.AuthGRPC == "" {
		return fmt.Errorf("--auth-grpc is required")
	}
	if Shop.PostgresDSN == "" {
		return fmt.Errorf("--postgres is required")
	}
	if Shop.KafkaAddr == "" {
		return fmt.Errorf("--kafka is required")
	}
	return nil
}

func LoadWorker() error {
	flag.StringVar(&Worker.PostgresDSN, "postgres", "", "postgres DSN")
	flag.StringVar(&Worker.KafkaAddr, "kafka", "", "kafka brokers (comma-separated host:port)")
	flag.StringVar(&Worker.OrdersTopic, "orders-topic", "orders.created", "kafka topic with new orders")
	flag.StringVar(&Worker.GroupID, "group-id", "shopx-workers", "kafka consumer group id")
	flag.IntVar(&Worker.Workers, "workers", 3, "number of concurrent workers")
	flag.Parse()

	if Worker.PostgresDSN == "" {
		return fmt.Errorf("--postgres is required")
	}
	if Worker.KafkaAddr == "" {
		return fmt.Errorf("--kafka is required")
	}
	if Worker.Workers <= 0 {
		return fmt.Errorf("--workers must be > 0")
	}
	return nil
}