package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RegisterRequest struct {
	Id       uint64
	Login    string
	Email    string
	Password string
}

type LoginRequest struct {
	LoginOrEmail string
	Password     string
}

type RefreshRequest struct {
	RefreshToken string
}

func (r *RefreshRequest) GetRefreshToken() string {
	if r == nil {
		return ""
	}

	return r.RefreshToken
}

type ValidateRequest struct {
	AccessToken string
}

type ValidateResponse struct {
	UserId uint64
}

type DeleteRequest struct {
	Id uint64
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthClient interface {
	Register(context.Context, *RegisterRequest, ...grpc.CallOption) (*emptypb.Empty, error)
	Login(context.Context, *LoginRequest, ...grpc.CallOption) (*TokenPair, error)
	Refresh(context.Context, *RefreshRequest, ...grpc.CallOption) (*TokenPair, error)
	Validate(context.Context, *ValidateRequest, ...grpc.CallOption) (*ValidateResponse, error)
	Logout(context.Context, *RefreshRequest, ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteUser(context.Context, *DeleteRequest, ...grpc.CallOption) (*emptypb.Empty, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc: cc}
}

func (c *authClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/auth.Auth/Register", in, out, opts...)
	return out, err
}

func (c *authClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*TokenPair, error) {
	out := new(TokenPair)
	err := c.cc.Invoke(ctx, "/auth.Auth/Login", in, out, opts...)
	return out, err
}

func (c *authClient) Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*TokenPair, error) {
	out := new(TokenPair)
	err := c.cc.Invoke(ctx, "/auth.Auth/Refresh", in, out, opts...)
	return out, err
}

func (c *authClient) Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, "/auth.Auth/Validate", in, out, opts...)
	return out, err
}

func (c *authClient) Logout(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/auth.Auth/Logout", in, out, opts...)
	return out, err
}

func (c *authClient) DeleteUser(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/auth.Auth/DeleteUser", in, out, opts...)
	return out, err
}
