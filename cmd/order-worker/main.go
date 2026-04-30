package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/EgorTarasov/shopx/internal/adapter/event/kafkasub"
	"github.com/EgorTarasov/shopx/internal/adapter/processor"
	"github.com/EgorTarasov/shopx/internal/config"
	"github.com/EgorTarasov/shopx/internal/configure"
	"github.com/EgorTarasov/shopx/internal/event"
	"github.com/EgorTarasov/shopx/internal/repository"
	"github.com/EgorTarasov/shopx/internal/usecase/worker"
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

	if err := config.LoadWorker(); err != nil {
		logger.Error("invalid configuration", "err", err)
		return exitError
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := serve(ctx, logger); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("order-worker exited with error", "err", err)
		return exitError
	}
	logger.Info("order-worker stopped")
	return exitSuccess
}

func serve(ctx context.Context, logger *slog.Logger) error {
	logger.Info("starting order-worker",
		"orders-topic", config.Worker.OrdersTopic,
		"group-id", config.Worker.GroupID,
		"workers", config.Worker.Workers,
	)

	db, _ := configure.NewPostgres(config.Worker.PostgresDSN)
	defer db.Close()

	tr := configure.NewKafkaTransport(config.Worker.KafkaAddr, config.Worker.GroupID)
	defer tr.Close()

	orders := repository.NewOrderRepo(db)
	proc := processor.NewStub(500*time.Millisecond, 0.1)
	svc := worker.NewService(orders, proc, logger)

	consumer := kafkasub.New(tr, config.Worker.OrdersTopic)

	jobs := make(chan event.OrderCreatedEvent, config.Worker.Workers)

	var wg sync.WaitGroup
	for i := 0; i < config.Worker.Workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			workerLog := logger.With("worker", id)
			for evt := range jobs {
				if err := svc.HandleCreated(ctx, worker.HandleCreatedInput{OrderID: evt.ID}); err != nil {
					workerLog.Error("handle order failed", "order_id", evt.ID, "err", err)
				}
			}
		}(i)
	}

	consumeErr := make(chan error, 1)
	go func() {
		consumeErr <- consumer.Run(ctx, func(ctx context.Context, e event.OrderCreatedEvent) error {
			select {
			case jobs <- e:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}()

	var runErr error
	select {
	case <-ctx.Done():
		runErr = ctx.Err()
	case err := <-consumeErr:
		if err != nil && !errors.Is(err, context.Canceled) {
			runErr = fmt.Errorf("kafka consume: %w", err)
		}
	}

	close(jobs)
	wg.Wait()
	return runErr
}
