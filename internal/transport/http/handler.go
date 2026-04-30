package httptransport

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
	api "github.com/EgorTarasov/shopx/internal/generated/api"
	"github.com/EgorTarasov/shopx/internal/usecase/shop"
)

type Handler struct {
	svc *shop.Service
}

func NewHandler(svc *shop.Service) *Handler { return &Handler{svc: svc} }

var _ api.StrictServerInterface = (*Handler)(nil)

func (h *Handler) requireUser(ctx context.Context) (shop.MeOutput, error) {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return shop.MeOutput{}, domain.ErrSessionNotFound
	}
	return h.svc.Me(ctx, shop.MeInput{Token: token})
}
