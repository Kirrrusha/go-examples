package account

import (
	"context"

	pagination "gateway/pkg/pagination/go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateUserRequest struct {
	User *CreateUser
}

type CreateUser struct {
	Login      string
	Phone      string
	FirstName  string
	LastName   string
	MiddleName string
	Email      string
	Age        uint32
}

type UpdateUserRequest struct {
	UserId uint64
	User   *UpdateUser
}

type UpdateUser struct {
	Phone      string
	FirstName  string
	LastName   string
	MiddleName string
	Email      string
	Age        uint32
}

type GetUserRequest struct {
	UserId uint64
}

type GetUserResponse struct {
	User *User
}

type GetUsersRequest struct {
	Pagination *pagination.Pagination
}

type GetUsersResponse struct {
	Users      []*User
	Pagination *pagination.Pagination
}

type DeleteUserRequest struct {
	UserId uint64
}

type User struct {
	Id         uint64
	Login      string
	Phone      string
	FirstName  string
	LastName   string
	MiddleName string
	Email      string
	Age        uint32
	CreatedAt  *timestamppb.Timestamp
	UpdatedAt  *timestamppb.Timestamp
}

type AccountClient interface {
	CreateUser(context.Context, *CreateUserRequest, ...grpc.CallOption) (*GetUserResponse, error)
	GetUser(context.Context, *GetUserRequest, ...grpc.CallOption) (*GetUserResponse, error)
	GetUsers(context.Context, *GetUsersRequest, ...grpc.CallOption) (*GetUsersResponse, error)
	DeleteUser(context.Context, *DeleteUserRequest, ...grpc.CallOption) (*emptypb.Empty, error)
	UpdateUser(context.Context, *UpdateUserRequest, ...grpc.CallOption) (*emptypb.Empty, error)
}

type accountClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountClient(cc grpc.ClientConnInterface) AccountClient {
	return &accountClient{cc: cc}
}

func (c *accountClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, "/account.Account/CreateUser", in, out, opts...)
	return out, err
}

func (c *accountClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, "/account.Account/GetUser", in, out, opts...)
	return out, err
}

func (c *accountClient) GetUsers(ctx context.Context, in *GetUsersRequest, opts ...grpc.CallOption) (*GetUsersResponse, error) {
	out := new(GetUsersResponse)
	err := c.cc.Invoke(ctx, "/account.Account/GetUsers", in, out, opts...)
	return out, err
}

func (c *accountClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/account.Account/DeleteUser", in, out, opts...)
	return out, err
}

func (c *accountClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/account.Account/UpdateUser", in, out, opts...)
	return out, err
}
