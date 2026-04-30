package httptransport

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
	api "github.com/EgorTarasov/shopx/internal/generated/api"
)

func (h *Handler) CreateOrder(ctx context.Context, req api.CreateOrderRequestObject) (api.CreateOrderResponseObject, error) {
	user, err := h.requireUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.Body == nil {
		return nil, domain.ErrInvalidInput
	}
	out, err := h.svc.CreateOrder(ctx, toCreateOrderInput(user.ID, *req.Body))
	if err != nil {
		return nil, err
	}
	return api.CreateOrder202JSONResponse(api.CreateOrderResponse{
		Id:     out.ID.String(),
		Status: api.OrderStatus(out.Status),
	}), nil
}
