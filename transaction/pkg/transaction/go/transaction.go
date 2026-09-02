package transaction

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	pagination "transaction/pkg/pagination/go"
)

type DepositRequest struct {
	UserId uint64
	Amount int64
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

type WithdrawRequest struct {
	UserId uint64
	Amount int64
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

type TransferRequest struct {
	UserId    uint64
	Amount    int64
	Recipient uint64
}

func (r *TransferRequest) GetUserId() uint64 {
	if r == nil {
		return 0
	}

	return r.UserId
}

func (r *TransferRequest) GetAmount() int64 {
	if r == nil {
		return 0
	}

	return r.Amount
}

func (r *TransferRequest) GetRecipient() uint64 {
	if r == nil {
		return 0
	}

	return r.Recipient
}

type GetTransactionsRequest struct {
	UserId     *uint64
	Type       *string
	Status     *string
	DateFrom   *timestamppb.Timestamp
	DateTo     *timestamppb.Timestamp
	Pagination *pagination.Pagination
}

func (r *GetTransactionsRequest) GetUserId() uint64 {
	if r == nil || r.UserId == nil {
		return 0
	}

	return *r.UserId
}

func (r *GetTransactionsRequest) GetType() string {
	if r == nil || r.Type == nil {
		return ""
	}

	return *r.Type
}

func (r *GetTransactionsRequest) GetStatus() string {
	if r == nil || r.Status == nil {
		return ""
	}

	return *r.Status
}

func (r *GetTransactionsRequest) GetDateFrom() *timestamppb.Timestamp {
	if r == nil {
		return nil
	}

	return r.DateFrom
}

func (r *GetTransactionsRequest) GetDateTo() *timestamppb.Timestamp {
	if r == nil {
		return nil
	}

	return r.DateTo
}

func (r *GetTransactionsRequest) GetPagination() *pagination.Pagination {
	if r == nil {
		return nil
	}

	return r.Pagination
}

type GetTransactionsResponse struct {
	Transactions []*TransactionDetails
}

type Transaction struct {
	Id        uint64
	UserId    uint64
	Amount    int64
	Type      string
	Status    string
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
}

type TransactionEntry struct {
	Id            uint64
	TransactionId uint64
	AccountId     uint64
	Direction     string
	Amount        int64
	CreatedAt     *timestamppb.Timestamp
	UpdatedAt     *timestamppb.Timestamp
}

type TransactionDetails struct {
	Transaction *Transaction
	Entries     []*TransactionEntry
}

type TransactionServiceServer interface {
	Deposit(context.Context, *DepositRequest) (*emptypb.Empty, error)
	Withdraw(context.Context, *WithdrawRequest) (*emptypb.Empty, error)
	Transfer(context.Context, *TransferRequest) (*emptypb.Empty, error)
	GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsResponse, error)
}

type UnimplementedTransactionServiceServer struct{}

func (UnimplementedTransactionServiceServer) Deposit(context.Context, *DepositRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}

func (UnimplementedTransactionServiceServer) Withdraw(context.Context, *WithdrawRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}

func (UnimplementedTransactionServiceServer) Transfer(context.Context, *TransferRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transfer not implemented")
}

func (UnimplementedTransactionServiceServer) GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactions not implemented")
}

func RegisterTransactionServiceServer(s grpc.ServiceRegistrar, srv TransactionServiceServer) {
	s.RegisterService(&TransactionService_ServiceDesc, srv)
}

var TransactionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "transaction.TransactionService",
	HandlerType: (*TransactionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Deposit", Handler: depositHandler},
		{MethodName: "Withdraw", Handler: withdrawHandler},
		{MethodName: "Transfer", Handler: transferHandler},
		{MethodName: "GetTransactions", Handler: getTransactionsHandler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "transaction/transaction.proto",
}

func depositHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/transaction.TransactionService/Deposit"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).Deposit(ctx, req.(*DepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func withdrawHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/transaction.TransactionService/Withdraw"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).Withdraw(ctx, req.(*WithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func transferHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).Transfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/transaction.TransactionService/Transfer"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).Transfer(ctx, req.(*TransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func getTransactionsHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServiceServer).GetTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/transaction.TransactionService/GetTransactions"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServiceServer).GetTransactions(ctx, req.(*GetTransactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}
