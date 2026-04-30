package grpctransport

import (
	"errors"

	"github.com/EgorTarasov/shopx/internal/domain"
	authv1 "github.com/EgorTarasov/shopx/internal/generated/auth/v1"
	"github.com/EgorTarasov/shopx/internal/usecase/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	svc *auth.Service
}

func NewAuthServer(svc *auth.Service) *AuthServer { return &AuthServer{svc: svc} }

func mapErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, domain.ErrEmailTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrBadCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrSessionNotFound),
		errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
