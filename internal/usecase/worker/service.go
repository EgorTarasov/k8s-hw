package worker

import "log/slog"

type Service struct {
	orders    OrderRepo
	processor OrderProcessor
	logger    *slog.Logger
}

func NewService(orders OrderRepo, processor OrderProcessor, logger *slog.Logger) *Service {
	return &Service{
		orders:    orders,
		processor: processor,
		logger:    logger,
	}
}
