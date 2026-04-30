package configure

import (
	"context"
	"fmt"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/larek-tech/storage/postgres"
)

func NewPostgres(dsn string) (*postgres.DB, *manager.Manager) {
	ctx := context.Background()

	cfg, err := postgres.NewCfgFromDSN(dsn)
	if err != nil {
		panic(fmt.Errorf("failed to parse database config: %w", err))
	}

	db, trManager, err := postgres.New(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}
	return db, trManager
}
