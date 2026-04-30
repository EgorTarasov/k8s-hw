package httptransport

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
	api "github.com/EgorTarasov/shopx/internal/generated/api"
	"github.com/EgorTarasov/shopx/internal/usecase/shop"
)

func (h *Handler) GetOrder(ctx context.Context, req api.GetOrderRequestObject) (api.GetOrderResponseObject, error) {
	user, err := h.requireUser(ctx)
	if err != nil {
		return nil, err
	}
	out, err := h.svc.GetOrder(ctx, shop.GetOrderInput{
		UserID: user.ID,
		ID:     domain.OrderID(req.Id),
	})
	if err != nil {
		return nil, err
	}
	return api.GetOrder200JSONResponse(fromGetOrderOutput(out)), nil
}
