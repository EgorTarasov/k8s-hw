package worker

import (
	"context"
	"fmt"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type HandleCreatedInput struct {
	OrderID domain.OrderID
}

func (s *Service) HandleCreated(ctx context.Context, in HandleCreatedInput) error {
	o, err := s.orders.GetByID(ctx, in.OrderID)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	if !o.Status.CanTransitionTo(domain.StatusProcessing) {
		s.logger.Info("skip already-handled order", "order_id", o.ID, "status", o.Status)
		return nil
	}

	if err := s.orders.UpdateStatus(ctx, o.ID, domain.StatusProcessing); err != nil {
		return fmt.Errorf("set processing: %w", err)
	}

	if procErr := s.processor.Process(ctx, o); procErr != nil {
		s.logger.Warn("processor failed", "order_id", o.ID, "err", procErr)
		if err := s.orders.UpdateStatus(ctx, o.ID, domain.StatusFailed); err != nil {
			return fmt.Errorf("set failed: %w", err)
		}
		return procErr
	}

	if err := s.orders.UpdateStatus(ctx, o.ID, domain.StatusProcessed); err != nil {
		return fmt.Errorf("set processed: %w", err)
	}
	s.logger.Info("order processed", "order_id", o.ID)
	return nil
}
