package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/EgorTarasov/shopx/internal/configure"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	dsn := os.Getenv("MIGRATOR_POSTGRES_DSN")
	if dsn == "" {
		logger.Error("MIGRATOR_POSTGRES_DSN is required")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, _ := configure.NewPostgres(dsn)
	defer db.Close()

	logger.Info("applying migrations")
	if err := configure.Migrate(ctx, db); err != nil {
		logger.Error("migrate", "err", err)
		os.Exit(1)
	}
	logger.Info("migrations applied")
}
