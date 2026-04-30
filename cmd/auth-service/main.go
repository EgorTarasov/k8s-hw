package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EgorTarasov/shopx/internal/adapter/crypto"
	"github.com/EgorTarasov/shopx/internal/config"
	"github.com/EgorTarasov/shopx/internal/configure"
	authv1 "github.com/EgorTarasov/shopx/internal/generated/auth/v1"
	"github.com/EgorTarasov/shopx/internal/repository"
	grpctransport "github.com/EgorTarasov/shopx/internal/transport/grpc"
	"github.com/EgorTarasov/shopx/internal/usecase/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	if err := config.LoadAuth(); err != nil {
		logger.Error("invalid configuration", "err", err)
		return exitError
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := serve(ctx, logger); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("auth-service exited with error", "err", err)
		return exitError
	}
	logger.Info("auth-service stopped")
	return exitSuccess
}

func serve(ctx context.Context, logger *slog.Logger) error {
	logger.Info("starting auth-service", "listen", config.Auth.Listen)

	db, _ := configure.NewPostgres(config.Auth.PostgresDSN)
	defer db.Close()
	if err := configure.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	rds := configure.NewRedis(config.Auth.RedisDSN)
	defer rds.Close()

	users := repository.NewUserRepo(db)
	sessions := repository.NewSessionRepo(rds)
	hasher := crypto.NewBcryptHasher()
	ids := crypto.NewUUIDGen()

	svc := auth.NewService(users, sessions, hasher, ids, ids)
	authServer := grpctransport.NewAuthServer(svc)

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.Auth.Listen)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("auth gRPC server listening", "addr", config.Auth.Listen)
		errCh <- grpcServer.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutting down auth gRPC server")
		stopped := make(chan struct{})
		go func() {
			grpcServer.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
		case <-time.After(5 * time.Second):
			grpcServer.Stop()
		}
		return ctx.Err()
	case err := <-errCh:
		return fmt.Errorf("grpc serve: %w", err)
	}
}
