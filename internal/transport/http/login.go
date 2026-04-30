package httptransport

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
	api "github.com/EgorTarasov/shopx/internal/generated/api"
	"github.com/EgorTarasov/shopx/internal/usecase/shop"
)

func (h *Handler) LoginUser(ctx context.Context, req api.LoginUserRequestObject) (api.LoginUserResponseObject, error) {
	if req.Body == nil {
		return nil, domain.ErrInvalidInput
	}
	out, err := h.svc.Login(ctx, shop.LoginInput{
		Email:    string(req.Body.Email),
		Password: req.Body.Password,
	})
	if err != nil {
		return nil, err
	}
	return api.LoginUser200JSONResponse{SessionToken: out.Token.String()}, nil
}
