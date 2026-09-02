package account

import (
	"context"

	"google.golang.org/grpc"
)

type DepositRequest struct {
	UserId      uint64
	Amount      int64
	OperationId string
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

type TransferResponse struct {
	Status           string
	UserBalance      int64
	RecipientBalance int64
}

type GetBalanceRequest struct {
	UserId uint64
}

type GetBalanceResponse struct {
	Balance int64
}

type AccountClient interface {
	Deposit(context.Context, *DepositRequest, ...grpc.CallOption) (*DepositResponse, error)
	Withdraw(context.Context, *WithdrawRequest, ...grpc.CallOption) (*WithdrawResponse, error)
	Transfer(context.Context, *TransferRequest, ...grpc.CallOption) (*TransferResponse, error)
	GetBalance(context.Context, *GetBalanceRequest, ...grpc.CallOption) (*GetBalanceResponse, error)
}

type accountClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountClient(cc grpc.ClientConnInterface) AccountClient {
	return &accountClient{cc: cc}
}

func (c *accountClient) Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositResponse, error) {
	out := new(DepositResponse)
	err := c.cc.Invoke(ctx, "/account.Account/Deposit", in, out, opts...)
	return out, err
}

func (c *accountClient) Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawResponse, error) {
	out := new(WithdrawResponse)
	err := c.cc.Invoke(ctx, "/account.Account/Withdraw", in, out, opts...)
	return out, err
}

func (c *accountClient) Transfer(ctx context.Context, in *TransferRequest, opts ...grpc.CallOption) (*TransferResponse, error) {
	out := new(TransferResponse)
	err := c.cc.Invoke(ctx, "/account.Account/Transfer", in, out, opts...)
	return out, err
}

func (c *accountClient) GetBalance(ctx context.Context, in *GetBalanceRequest, opts ...grpc.CallOption) (*GetBalanceResponse, error) {
	out := new(GetBalanceResponse)
	err := c.cc.Invoke(ctx, "/account.Account/GetBalance", in, out, opts...)
	return out, err
}
