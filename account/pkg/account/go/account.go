package account

import (
	"context"

	pagination "account/pkg/pagination/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CreateUserRequest struct {
	User *CreateUser
}

func (r *CreateUserRequest) GetUser() *CreateUser {
	if r == nil {
		return nil
	}

	return r.User
}

type CreateUserResponse struct {
	User *User
}

type CreateUser struct {
	Login      string
	Phone      string
	FirstName  string
	LastName   string
	MiddleName string
	Email      string
	Age        uint32
	Balance    int64
	Password   string
}

type UpdateUserRequest struct {
	UserId uint64
	User   *UpdateUser
}

func (r *UpdateUserRequest) GetUserId() uint64 {
	if r == nil {
		return 0
	}

	return r.UserId
}

type UpdateUser struct {
	Phone      string
	FirstName  string
	LastName   string
	MiddleName string
	Email      string
	Age        uint32
	Balance    int64
}

type GetUserRequest struct {
	UserId uint64
}

func (r *GetUserRequest) GetUserId() uint64 {
	if r == nil {
		return 0
	}

	return r.UserId
}

type GetUserResponse struct {
	User *User
}

type GetUsersRequest struct {
	Pagination *pagination.Pagination
}

func (r *GetUsersRequest) GetPagination() *pagination.Pagination {
	if r == nil {
		return nil
	}

	return r.Pagination
}

type GetUsersResponse struct {
	Users      []*User
	Pagination *pagination.Pagination
}

type DeleteUserRequest struct {
	UserId uint64
}

func (r *DeleteUserRequest) GetUserId() uint64 {
	if r == nil {
		return 0
	}

	return r.UserId
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
	Balance    int64
	IsDeleted  bool
	CreatedAt  *timestamppb.Timestamp
	UpdatedAt  *timestamppb.Timestamp
}

type DepositRequest struct {
	UserId      uint64
	Amount      int64
	OperationId string
}

func (r *DepositRequest) GetUserId() uint64 {
	if r == nil {
		return 0
	}

	return r.UserId
}

func (r *DepositRequest) GetAmount() int64 {
	if r == nil {
		return 0
	}

	return r.Amount
}

type DepositResponse struct {
	Status  string
	Balance int64
}

type WithdrawRequest struct {
	UserId      uint64
	Amount      int64
	OperationId string
}

func (r *WithdrawRequest) GetUserId() uint64 {
	if r == nil {
		return 0
	}

	return r.UserId
}

func (r *WithdrawRequest) GetAmount() int64 {
	if r == nil {
		return 0
	}

	return r.Amount
}

type WithdrawResponse struct {
	Status  string
	Balance int64
}

type TransferRequest struct {
	UserId      uint64
	RecipientId uint64
	Amount      int64
	OperationId string
}

func (r *TransferRequest) GetUserId() uint64 {
	if r == nil {
		return 0
	}

	return r.UserId
}

func (r *TransferRequest) GetRecipientId() uint64 {
	if r == nil {
		return 0
	}

	return r.RecipientId
}

func (r *TransferRequest) GetAmount() int64 {
	if r == nil {
		return 0
	}

	return r.Amount
}

type TransferResponse struct {
	Status           string
	UserBalance      int64
	RecipientBalance int64
}

type GetBalanceRequest struct {
	UserId uint64
}

func (r *GetBalanceRequest) GetUserId() uint64 {
	if r == nil {
		return 0
	}

	return r.UserId
}

type GetBalanceResponse struct {
	Balance int64
}

type AccountServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*emptypb.Empty, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*emptypb.Empty, error)
	Deposit(context.Context, *DepositRequest) (*DepositResponse, error)
	Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error)
	Transfer(context.Context, *TransferRequest) (*TransferResponse, error)
	GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error)
}

type UnimplementedAccountServer struct{}

func (UnimplementedAccountServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}

func (UnimplementedAccountServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}

func (UnimplementedAccountServer) GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}

func (UnimplementedAccountServer) DeleteUser(context.Context, *DeleteUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}

func (UnimplementedAccountServer) UpdateUser(context.Context, *UpdateUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}

func (UnimplementedAccountServer) Deposit(context.Context, *DepositRequest) (*DepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}

func (UnimplementedAccountServer) Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}

func (UnimplementedAccountServer) Transfer(context.Context, *TransferRequest) (*TransferResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transfer not implemented")
}

func (UnimplementedAccountServer) GetBalance(context.Context, *GetBalanceRequest) (*GetBalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}

func RegisterAccountServer(s grpc.ServiceRegistrar, srv AccountServer) {
	s.RegisterService(&Account_ServiceDesc, srv)
}

var Account_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "account.Account",
	HandlerType: (*AccountServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateUser", Handler: createUserHandler},
		{MethodName: "GetUser", Handler: getUserHandler},
		{MethodName: "GetUsers", Handler: getUsersHandler},
		{MethodName: "DeleteUser", Handler: deleteUserHandler},
		{MethodName: "UpdateUser", Handler: updateUserHandler},
		{MethodName: "Deposit", Handler: depositHandler},
		{MethodName: "Withdraw", Handler: withdrawHandler},
		{MethodName: "Transfer", Handler: transferHandler},
		{MethodName: "GetBalance", Handler: getBalanceHandler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account/account.proto",
}

func createUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/account.Account/CreateUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func getUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/account.Account/GetUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func getUsersHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).GetUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/account.Account/GetUsers"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).GetUsers(ctx, req.(*GetUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func deleteUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/account.Account/DeleteUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func updateUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/account.Account/UpdateUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func depositHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/account.Account/Deposit"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).Deposit(ctx, req.(*DepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func withdrawHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/account.Account/Withdraw"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).Withdraw(ctx, req.(*WithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func transferHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).Transfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/account.Account/Transfer"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).Transfer(ctx, req.(*TransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func getBalanceHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/account.Account/GetBalance"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServer).GetBalance(ctx, req.(*GetBalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}
