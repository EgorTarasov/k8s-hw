package grpctransport

import (
	"context"

	authv1 "github.com/EgorTarasov/shopx/internal/generated/auth/v1"
	"github.com/EgorTarasov/shopx/internal/usecase/auth"
)

func (s *AuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	out, err := s.svc.Login(ctx, auth.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, mapErr(err)
	}
	return &authv1.LoginResponse{SessionToken: out.Token.String()}, nil
}
