package processor

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/EgorTarasov/shopx/internal/domain"
)

type Stub struct {
	delay     time.Duration
	failRatio float64
	rng       *rand.Rand
}

func NewStub(delay time.Duration, failRatio float64) *Stub {
	return &Stub{
		delay:     delay,
		failRatio: failRatio,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Stub) Process(ctx context.Context, _ domain.Order) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(s.delay):
	}
	if s.rng.Float64() < s.failRatio {
		return errors.New("simulated processing error")
	}
	return nil
}
