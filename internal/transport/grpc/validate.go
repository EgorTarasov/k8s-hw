package grpctransport

import (
	"context"

	"github.com/EgorTarasov/shopx/internal/domain"
	authv1 "github.com/EgorTarasov/shopx/internal/generated/auth/v1"
	"github.com/EgorTarasov/shopx/internal/usecase/auth"
)

func (s *AuthServer) Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	out, err := s.svc.Validate(ctx, auth.ValidateInput{Token: domain.SessionToken(req.GetSessionToken())})
	if err != nil {
		return nil, mapErr(err)
	}
	return &authv1.ValidateResponse{User: &authv1.User{
		Id:    out.ID.String(),
		Email: out.Email,
		Name:  out.Name,
	}}, nil
}
