package grpctransport

import (
	"context"

	authv1 "github.com/EgorTarasov/shopx/internal/generated/auth/v1"
	"github.com/EgorTarasov/shopx/internal/usecase/auth"
)

func (s *AuthServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	out, err := s.svc.Register(ctx, auth.RegisterInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Name:     req.GetName(),
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &authv1.RegisterResponse{User: &authv1.User{
		Id:    out.ID.String(),
		Email: out.Email,
		Name:  out.Name,
	}}, nil
}
