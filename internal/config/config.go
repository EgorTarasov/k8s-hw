package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
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

func envStr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func LoadAuth() error {
	flag.StringVar(&Auth.Listen, "listen", envStr("AUTH_LISTEN", "127.0.0.1:9090"), "gRPC listen address (env AUTH_LISTEN)")
	flag.StringVar(&Auth.PostgresDSN, "postgres", envStr("AUTH_POSTGRES_DSN", ""), "postgres DSN (env AUTH_POSTGRES_DSN)")
	flag.StringVar(&Auth.RedisDSN, "redis", envStr("AUTH_REDIS_ADDR", ""), "redis address host:port (env AUTH_REDIS_ADDR)")
	flag.Parse()

	if Auth.PostgresDSN == "" {
		return fmt.Errorf("postgres DSN is required (--postgres or AUTH_POSTGRES_DSN)")
	}
	if Auth.RedisDSN == "" {
		return fmt.Errorf("redis address is required (--redis or AUTH_REDIS_ADDR)")
	}
	return nil
}

func LoadShop() error {
	flag.StringVar(&Shop.Listen, "listen", envStr("SHOP_LISTEN", "127.0.0.1:8080"), "HTTP listen address (env SHOP_LISTEN)")
	flag.StringVar(&Shop.AuthGRPC, "auth-grpc", envStr("SHOP_AUTH_GRPC", ""), "auth-service gRPC address (env SHOP_AUTH_GRPC)")
	flag.StringVar(&Shop.PostgresDSN, "postgres", envStr("SHOP_POSTGRES_DSN", ""), "postgres DSN (env SHOP_POSTGRES_DSN)")
	flag.StringVar(&Shop.KafkaAddr, "kafka", envStr("SHOP_KAFKA_BROKERS", ""), "kafka brokers comma-separated (env SHOP_KAFKA_BROKERS)")
	flag.StringVar(&Shop.OrdersTopic, "orders-topic", envStr("SHOP_ORDERS_TOPIC", "orders.created"), "kafka topic for created orders (env SHOP_ORDERS_TOPIC)")
	flag.Parse()

	if Shop.AuthGRPC == "" {
		return fmt.Errorf("auth-grpc address is required (--auth-grpc or SHOP_AUTH_GRPC)")
	}
	if Shop.PostgresDSN == "" {
		return fmt.Errorf("postgres DSN is required (--postgres or SHOP_POSTGRES_DSN)")
	}
	if Shop.KafkaAddr == "" {
		return fmt.Errorf("kafka brokers are required (--kafka or SHOP_KAFKA_BROKERS)")
	}
	return nil
}

func LoadWorker() error {
	flag.StringVar(&Worker.PostgresDSN, "postgres", envStr("WORKER_POSTGRES_DSN", ""), "postgres DSN (env WORKER_POSTGRES_DSN)")
	flag.StringVar(&Worker.KafkaAddr, "kafka", envStr("WORKER_KAFKA_BROKERS", ""), "kafka brokers comma-separated (env WORKER_KAFKA_BROKERS)")
	flag.StringVar(&Worker.OrdersTopic, "orders-topic", envStr("WORKER_ORDERS_TOPIC", "orders.created"), "kafka topic with new orders (env WORKER_ORDERS_TOPIC)")
	flag.StringVar(&Worker.GroupID, "group-id", envStr("WORKER_GROUP_ID", "shopx-workers"), "kafka consumer group id (env WORKER_GROUP_ID)")
	flag.IntVar(&Worker.Workers, "workers", envInt("WORKER_CONCURRENCY", 3), "number of concurrent workers (env WORKER_CONCURRENCY)")
	flag.Parse()

	if Worker.PostgresDSN == "" {
		return fmt.Errorf("postgres DSN is required (--postgres or WORKER_POSTGRES_DSN)")
	}
	if Worker.KafkaAddr == "" {
		return fmt.Errorf("kafka brokers are required (--kafka or WORKER_KAFKA_BROKERS)")
	}
	if Worker.Workers <= 0 {
		return fmt.Errorf("workers must be > 0 (--workers or WORKER_CONCURRENCY)")
	}
	return nil
}
