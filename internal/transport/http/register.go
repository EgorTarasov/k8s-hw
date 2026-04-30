package httptransport

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
	api "github.com/EgorTarasov/shopx/internal/generated/api"
	"github.com/EgorTarasov/shopx/internal/usecase/shop"
)

func (h *Handler) RegisterUser(ctx context.Context, req api.RegisterUserRequestObject) (api.RegisterUserResponseObject, error) {
	if req.Body == nil {
		return nil, domain.ErrInvalidInput
	}
	out, err := h.svc.Register(ctx, shop.RegisterInput{
		Email:    string(req.Body.Email),
		Password: req.Body.Password,
		Name:     req.Body.Name,
	})
	if err != nil {
		return nil, err
	}
	return api.RegisterUser201JSONResponse(api.User{
		Id:    out.ID.String(),
		Email: openapiEmail(out.Email),
		Name:  out.Name,
	}), nil
}
