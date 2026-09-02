package server

import (
	"context"

	"account/internal/mapper"
	"account/internal/model"
	accountpb "account/pkg/account/go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	accountpb.UnimplementedAccountServer
	accountService AccountService
	logger         *zerolog.Logger
}

func New(accountService AccountService, logger *zerolog.Logger) *Server {
	return &Server{accountService: accountService, logger: logger}
}

type AccountService interface {
	CreateUser(context.Context, model.CreateUser) (model.User, error)
	GetUser(context.Context, uint64) (model.User, error)
	GetUsers(context.Context, int, int) ([]model.User, error)
	DeleteUser(context.Context, uint64) error
	UpdateUser(context.Context, uint64, model.UpdateUser) error
	Deposit(context.Context, uint64, int64) (int64, error)
	Withdraw(context.Context, uint64, int64) (int64, error)
	Transfer(context.Context, uint64, uint64, int64) (int64, int64, error)
	GetBalance(context.Context, uint64) (int64, error)
}

func (s *Server) CreateUser(ctx context.Context, req *accountpb.CreateUserRequest) (*accountpb.CreateUserResponse, error) {
	user := mapper.PbToUserCreate(req.GetUser())
	createdUser, err := s.accountService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &accountpb.CreateUserResponse{User: mapper.UserToPb(createdUser)}, nil
}

func (s *Server) GetUser(ctx context.Context, req *accountpb.GetUserRequest) (*accountpb.GetUserResponse, error) {
	res, err := s.accountService.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &accountpb.GetUserResponse{User: mapper.UserToPb(res)}, nil
}

func (s *Server) GetUsers(ctx context.Context, req *accountpb.GetUsersRequest) (*accountpb.GetUsersResponse, error) {
	pg := req.GetPagination()
	res, err := s.accountService.GetUsers(ctx, int(pg.GetLimit()), int(pg.GetOffset()))
	if err != nil {
		return nil, err
	}

	return &accountpb.GetUsersResponse{
		Users:      mapper.UsersToPbs(res),
		Pagination: pg,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *accountpb.UpdateUserRequest) (*emptypb.Empty, error) {
	user := mapper.PbToUserUpdate(req.User)
	if err := s.accountService.UpdateUser(ctx, req.GetUserId(), user); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *accountpb.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := s.accountService.DeleteUser(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) Deposit(ctx context.Context, req *accountpb.DepositRequest) (*accountpb.DepositResponse, error) {
	balance, err := s.accountService.Deposit(ctx, req.GetUserId(), req.GetAmount())
	if err != nil {
		return nil, err
	}

	return &accountpb.DepositResponse{Status: "completed", Balance: balance}, nil
}

func (s *Server) Withdraw(ctx context.Context, req *accountpb.WithdrawRequest) (*accountpb.WithdrawResponse, error) {
	balance, err := s.accountService.Withdraw(ctx, req.GetUserId(), req.GetAmount())
	if err != nil {
		return nil, err
	}

	return &accountpb.WithdrawResponse{Status: "completed", Balance: balance}, nil
}

func (s *Server) Transfer(ctx context.Context, req *accountpb.TransferRequest) (*accountpb.TransferResponse, error) {
	userBalance, recipientBalance, err := s.accountService.Transfer(ctx, req.GetUserId(), req.GetRecipientId(), req.GetAmount())
	if err != nil {
		return nil, err
	}

	return &accountpb.TransferResponse{
		Status:           "completed",
		UserBalance:      userBalance,
		RecipientBalance: recipientBalance,
	}, nil
}

func (s *Server) GetBalance(ctx context.Context, req *accountpb.GetBalanceRequest) (*accountpb.GetBalanceResponse, error) {
	balance, err := s.accountService.GetBalance(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &accountpb.GetBalanceResponse{Balance: balance}, nil
}
