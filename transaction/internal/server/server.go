package server

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/emptypb"
	"transaction/internal/mapper"
	"transaction/internal/model"
	transactionpb "transaction/pkg/transaction/go"
)

type Server struct {
	transactionpb.UnimplementedTransactionServiceServer
	transactionService TransactionService
	logger             *zerolog.Logger
}

func New(transactionService TransactionService, logger *zerolog.Logger) *Server {
	return &Server{transactionService: transactionService, logger: logger}
}

type TransactionService interface {
	GetTransactionsWithDetails(context.Context, model.GetTransactionsParams) ([]model.TransactionDetails, error)
	Deposit(context.Context, uint64, int64) (model.TransactionDetails, error)
	Withdraw(context.Context, uint64, int64) (model.TransactionDetails, error)
	Transfer(context.Context, uint64, uint64, int64) (model.TransactionDetails, error)
}

func (s *Server) Deposit(ctx context.Context, req *transactionpb.DepositRequest) (*emptypb.Empty, error) {
	if _, err := s.transactionService.Deposit(ctx, req.GetUserId(), req.GetAmount()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) Withdraw(ctx context.Context, req *transactionpb.WithdrawRequest) (*emptypb.Empty, error) {
	if _, err := s.transactionService.Withdraw(ctx, req.GetUserId(), req.GetAmount()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) Transfer(ctx context.Context, req *transactionpb.TransferRequest) (*emptypb.Empty, error) {
	if _, err := s.transactionService.Transfer(ctx, req.GetUserId(), req.GetRecipient(), req.GetAmount()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetTransactions(ctx context.Context, req *transactionpb.GetTransactionsRequest) (*transactionpb.GetTransactionsResponse, error) {
	params := mapper.GetTransactionsRequestToParams(req)
	res, err := s.transactionService.GetTransactionsWithDetails(ctx, params)
	if err != nil {
		return nil, err
	}

	return &transactionpb.GetTransactionsResponse{
		Transactions: mapper.TransactionDetailsListToPb(res),
	}, nil
}
