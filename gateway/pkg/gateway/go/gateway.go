package gateway

import (
	"context"

	account "gateway/pkg/account/go"
	auth "gateway/pkg/auth/go"
	pagination "gateway/pkg/pagination/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RegisterRequest struct {
	User     *account.User
	Password string
}

func (r *RegisterRequest) GetUser() *account.User {
	if r == nil {
		return nil
	}

	return r.User
}

func (r *RegisterRequest) GetPassword() string {
	if r == nil {
		return ""
	}

	return r.Password
}

type RegisterResponse struct {
	User   *account.User
	Tokens *auth.TokenPair
}

type LoginRequest struct {
	LoginOrEmail string
	Password     string
}

func (r *LoginRequest) GetLoginOrEmail() string {
	if r == nil {
		return ""
	}

	return r.LoginOrEmail
}

func (r *LoginRequest) GetPassword() string {
	if r == nil {
		return ""
	}

	return r.Password
}

type LoginResponse struct {
	User   *account.User
	Tokens *auth.TokenPair
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

type RefreshResponse struct {
	TokenPair *auth.TokenPair
}

type LogoutRequest struct {
	RefreshToken string
}

func (r *LogoutRequest) GetRefreshToken() string {
	if r == nil {
		return ""
	}

	return r.RefreshToken
}

type ValidateTokenRequest struct {
	AccessToken string
}

func (r *ValidateTokenRequest) GetAccessToken() string {
	if r == nil {
		return ""
	}

	return r.AccessToken
}

type ValidateTokenResponse struct {
	UserId  uint64
	IsValid bool
}

type CreateUserRequest struct {
	User *account.CreateUser
}

func (r *CreateUserRequest) GetUser() *account.CreateUser {
	if r == nil {
		return nil
	}

	return r.User
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
	User *account.User
}

type GetCurrentUserResponse struct {
	User *account.User
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
	Users      []*account.User
	Pagination *pagination.Pagination
}

type UpdateUserRequest struct {
	UserId uint64
	User   *account.UpdateUser
}

func (r *UpdateUserRequest) GetUserId() uint64 {
	if r == nil {
		return 0
	}

	return r.UserId
}

func (r *UpdateUserRequest) GetUser() *account.UpdateUser {
	if r == nil {
		return nil
	}

	return r.User
}

type UpdateCurrentUserRequest struct {
	User *account.UpdateUser
}

func (r *UpdateCurrentUserRequest) GetUser() *account.UpdateUser {
	if r == nil {
		return nil
	}

	return r.User
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

type GatewayServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error)
	Logout(context.Context, *LogoutRequest) (*emptypb.Empty, error)
	ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error)
	CreateUser(context.Context, *CreateUserRequest) (*emptypb.Empty, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	GetCurrentUser(context.Context, *emptypb.Empty) (*GetCurrentUserResponse, error)
	GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error)
	UpdateUser(context.Context, *UpdateUserRequest) (*emptypb.Empty, error)
	UpdateCurrentUser(context.Context, *UpdateCurrentUserRequest) (*emptypb.Empty, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*emptypb.Empty, error)
	DeleteCurrentUser(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
}

type UnimplementedGatewayServer struct{}

func (UnimplementedGatewayServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}

func (UnimplementedGatewayServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}

func (UnimplementedGatewayServer) Refresh(context.Context, *RefreshRequest) (*RefreshResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Refresh not implemented")
}

func (UnimplementedGatewayServer) Logout(context.Context, *LogoutRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}

func (UnimplementedGatewayServer) ValidateToken(context.Context, *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}

func (UnimplementedGatewayServer) CreateUser(context.Context, *CreateUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}

func (UnimplementedGatewayServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}

func (UnimplementedGatewayServer) GetCurrentUser(context.Context, *emptypb.Empty) (*GetCurrentUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrentUser not implemented")
}

func (UnimplementedGatewayServer) GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}

func (UnimplementedGatewayServer) UpdateUser(context.Context, *UpdateUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}

func (UnimplementedGatewayServer) UpdateCurrentUser(context.Context, *UpdateCurrentUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCurrentUser not implemented")
}

func (UnimplementedGatewayServer) DeleteUser(context.Context, *DeleteUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}

func (UnimplementedGatewayServer) DeleteCurrentUser(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCurrentUser not implemented")
}

func RegisterGatewayServer(s grpc.ServiceRegistrar, srv GatewayServer) {
	s.RegisterService(&Gateway_ServiceDesc, srv)
}

var Gateway_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.Gateway",
	HandlerType: (*GatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Register", Handler: registerHandler},
		{MethodName: "Login", Handler: loginHandler},
		{MethodName: "Refresh", Handler: refreshHandler},
		{MethodName: "Logout", Handler: logoutHandler},
		{MethodName: "ValidateToken", Handler: validateTokenHandler},
		{MethodName: "CreateUser", Handler: createUserHandler},
		{MethodName: "GetUser", Handler: getUserHandler},
		{MethodName: "GetCurrentUser", Handler: getCurrentUserHandler},
		{MethodName: "GetUsers", Handler: getUsersHandler},
		{MethodName: "UpdateUser", Handler: updateUserHandler},
		{MethodName: "UpdateCurrentUser", Handler: updateCurrentUserHandler},
		{MethodName: "DeleteUser", Handler: deleteUserHandler},
		{MethodName: "DeleteCurrentUser", Handler: deleteCurrentUserHandler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gateway/gateway.proto",
}

func registerHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/Register"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func loginHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/Login"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func refreshHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).Refresh(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/Refresh"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).Refresh(ctx, req.(*RefreshRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func logoutHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/Logout"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func validateTokenHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).ValidateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/ValidateToken"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).ValidateToken(ctx, req.(*ValidateTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func createUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/CreateUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func getUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/GetUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func getCurrentUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).GetCurrentUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/GetCurrentUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).GetCurrentUser(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func getUsersHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).GetUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/GetUsers"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).GetUsers(ctx, req.(*GetUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func updateUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/UpdateUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func updateCurrentUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCurrentUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).UpdateCurrentUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/UpdateCurrentUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).UpdateCurrentUser(ctx, req.(*UpdateCurrentUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func deleteUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/DeleteUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func deleteCurrentUserHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServer).DeleteCurrentUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/gateway.Gateway/DeleteCurrentUser"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServer).DeleteCurrentUser(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}
