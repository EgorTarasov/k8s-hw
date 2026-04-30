package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EgorTarasov/shopx/internal/adapter/auth/grpcclient"
	"github.com/EgorTarasov/shopx/internal/adapter/crypto"
	"github.com/EgorTarasov/shopx/internal/adapter/event/kafkapub"
	"github.com/EgorTarasov/shopx/internal/config"
	"github.com/EgorTarasov/shopx/internal/configure"
	"github.com/EgorTarasov/shopx/internal/repository"
	httptransport "github.com/EgorTarasov/shopx/internal/transport/http"
	"github.com/EgorTarasov/shopx/internal/usecase/shop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	exitSuccess = 0
	exitError   = 1
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	if err := config.LoadShop(); err != nil {
		logger.Error("invalid configuration", "err", err)
		return exitError
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := serve(ctx, logger); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("shop-backend exited with error", "err", err)
		return exitError
	}
	logger.Info("shop-backend stopped")
	return exitSuccess
}

func serve(ctx context.Context, logger *slog.Logger) error {
	logger.Info("starting shop-backend",
		"listen", config.Shop.Listen,
		"auth-grpc", config.Shop.AuthGRPC,
		"orders-topic", config.Shop.OrdersTopic,
	)

	db, _ := configure.NewPostgres(config.Shop.PostgresDSN)
	defer db.Close()

	authConn, err := grpc.NewClient(config.Shop.AuthGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("dial auth-service: %w", err)
	}
	defer authConn.Close()

	tr := configure.NewKafkaTransport(config.Shop.KafkaAddr, "")
	defer tr.Close()

	authClient := grpcclient.New(authConn)
	orders := repository.NewOrderRepo(db)
	publisher := kafkapub.New(tr, config.Shop.OrdersTopic)
	ids := crypto.NewUUIDGen()

	svc := shop.NewService(authClient, orders, publisher, ids, logger)
	handler := httptransport.NewHandler(svc)

	server := &http.Server{
		Addr:              config.Shop.Listen,
		Handler:           httptransport.Router(handler, logger),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("HTTP server listening", "addr", config.Shop.Listen)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("HTTP shutdown failed", "err", err)
		}
		return ctx.Err()
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("http serve: %w", err)
		}
		return nil
	}
}
