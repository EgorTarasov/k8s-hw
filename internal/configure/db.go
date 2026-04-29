package configure

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/larek-tech/storage/postgres"
)

func NewPostgres(dsn string) (*postgres.DB, *manager.Manager) {
	ctx := context.Background()

	cfg, err := parseDSN(dsn)
	if err != nil {
		panic(err)
	}

	db, trManager, err := postgres.New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return db, trManager
}

func parseDSN(dsn string) (postgres.Cfg, error) {
	parsed, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return postgres.Cfg{}, fmt.Errorf("parse dsn: %w", err)
	}
	return postgres.Cfg{
		User:     parsed.User,
		Password: parsed.Password,
		Host:     parsed.Host,
		Port:     int(parsed.Port),
		DB:       parsed.Database,
	}, nil
}
