package account

import (
	"context"
	"fmt"

	accountpb "transaction/pkg/account/go"
)

type Service struct {
	client accountpb.AccountClient
}

func New(client accountpb.AccountClient) *Service {
	return &Service{client: client}
}

func (s *Service) Deposit(ctx context.Context, userID uint64, amount int64, operationID string) error {
	if _, err := s.client.Deposit(ctx, &accountpb.DepositRequest{
		UserId:      userID,
		Amount:      amount,
		OperationId: operationID,
	}); err != nil {
		return fmt.Errorf("failed to deposit account balance: %w", err)
	}

	return nil
}

func (s *Service) Withdraw(ctx context.Context, userID uint64, amount int64, operationID string) error {
	if _, err := s.client.Withdraw(ctx, &accountpb.WithdrawRequest{
		UserId:      userID,
		Amount:      amount,
		OperationId: operationID,
	}); err != nil {
		return fmt.Errorf("failed to withdraw account balance: %w", err)
	}

	return nil
}

func (s *Service) Transfer(ctx context.Context, userID uint64, recipientID uint64, amount int64, operationID string) error {
	if _, err := s.client.Transfer(ctx, &accountpb.TransferRequest{
		UserId:      userID,
		RecipientId: recipientID,
		Amount:      amount,
		OperationId: operationID,
	}); err != nil {
		return fmt.Errorf("failed to transfer account balance: %w", err)
	}

	return nil
}
