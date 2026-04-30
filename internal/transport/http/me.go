package httptransport

import (
	"context"

	api "github.com/EgorTarasov/shopx/internal/generated/api"
)

func (h *Handler) GetMe(ctx context.Context, _ api.GetMeRequestObject) (api.GetMeResponseObject, error) {
	user, err := h.requireUser(ctx)
	if err != nil {
		return nil, err
	}
	return api.GetMe200JSONResponse(api.User{
		Id:    user.ID.String(),
		Email: openapiEmail(user.Email),
		Name:  user.Name,
	}), nil
}
