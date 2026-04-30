package grpcclient

import (
	"context"
	"errors"
	"fmt"

	authv1 "github.com/EgorTarasov/shopx/internal/generated/auth/v1"
	"github.com/EgorTarasov/shopx/internal/domain"
	"github.com/EgorTarasov/shopx/internal/usecase/shop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn *grpc.ClientConn
	api  authv1.AuthServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{conn: conn, api: authv1.NewAuthServiceClient(conn)}
}

func (c *Client) Register(ctx context.Context, in shop.RegisterInput) (shop.RegisterOutput, error) {
	resp, err := c.api.Register(ctx, &authv1.RegisterRequest{
		Email:    in.Email,
		Password: in.Password,
		Name:     in.Name,
	})
	if err != nil {
		return shop.RegisterOutput{}, mapErr(err)
	}
	if resp.User == nil {
		return shop.RegisterOutput{}, fmt.Errorf("auth-service: empty user in response")
	}
	return shop.RegisterOutput{
		ID:    domain.UserID(resp.User.Id),
		Email: resp.User.Email,
		Name:  resp.User.Name,
	}, nil
}

func (c *Client) Login(ctx context.Context, in shop.LoginInput) (shop.LoginOutput, error) {
	resp, err := c.api.Login(ctx, &authv1.LoginRequest{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		return shop.LoginOutput{}, mapErr(err)
	}
	return shop.LoginOutput{Token: domain.SessionToken(resp.SessionToken)}, nil
}

func (c *Client) Validate(ctx context.Context, token domain.SessionToken) (shop.MeOutput, error) {
	resp, err := c.api.Validate(ctx, &authv1.ValidateRequest{SessionToken: token.String()})
	if err != nil {
		return shop.MeOutput{}, mapErr(err)
	}
	if resp.User == nil {
		return shop.MeOutput{}, fmt.Errorf("auth-service: empty user in response")
	}
	return shop.MeOutput{
		ID:    domain.UserID(resp.User.Id),
		Email: resp.User.Email,
		Name:  resp.User.Name,
	}, nil
}

func mapErr(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case codes.AlreadyExists:
		return domain.ErrEmailTaken
	case codes.Unauthenticated:
		return domain.ErrBadCredentials
	case codes.NotFound:
		return domain.ErrSessionNotFound
	case codes.InvalidArgument:
		return errors.Join(domain.ErrInvalidInput, errors.New(st.Message()))
	default:
		return err
	}
}
