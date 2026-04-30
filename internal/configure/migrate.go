package configure

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/larek-tech/storage/postgres"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/shopx/*.sql
var migrationsFS embed.FS

func Migrate(ctx context.Context, db *postgres.DB) error {
	sqlDB := stdlib.OpenDBFromPool(db.GetPool())
	defer sqlDB.Close()

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	if err := goose.UpContext(ctx, sqlDB, "migrations/shopx"); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
