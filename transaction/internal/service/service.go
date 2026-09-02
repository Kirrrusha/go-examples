package service

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"transaction/internal/model"
)

type TransactionService struct {
	repo           Repository
	accountService AccountService
	logger         *zerolog.Logger
}

func New(repo Repository, accountService AccountService, logger *zerolog.Logger) *TransactionService {
	return &TransactionService{repo: repo, accountService: accountService, logger: logger}
}

type Repository interface {
	GetTransactions(context.Context, model.GetTransactionsParams) ([]model.Transaction, error)
	GetTransactionDetails(context.Context, uint64) (model.TransactionDetails, error)
	Deposit(context.Context, model.DepositParams) (model.TransactionDetails, error)
	Withdraw(context.Context, model.WithdrawParams) (model.TransactionDetails, error)
	Transfer(context.Context, model.TransferParams) (model.TransactionDetails, error)
}

type AccountService interface {
	Deposit(context.Context, uint64, int64, string) error
	Withdraw(context.Context, uint64, int64, string) error
	Transfer(context.Context, uint64, uint64, int64, string) error
}

func (s *TransactionService) GetTransactionsWithDetails(ctx context.Context, params model.GetTransactionsParams) ([]model.TransactionDetails, error) {
	transactions, err := s.repo.GetTransactions(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get transactions")
		return nil, err
	}

	transactionDetails := make([]model.TransactionDetails, 0, len(transactions))
	for _, transaction := range transactions {
		details, err := s.repo.GetTransactionDetails(ctx, transaction.ID)
		if err != nil {
			s.logger.Error().Err(err).Uint64("transaction_id", transaction.ID).Msg("failed to get transaction details")
			continue
		}

		transactionDetails = append(transactionDetails, details)
	}

	return transactionDetails, nil
}

func (s *TransactionService) Deposit(ctx context.Context, userID uint64, amount int64) (model.TransactionDetails, error) {
	params := model.DepositParams{UserID: userID, Amount: amount}
	details, err := s.repo.Deposit(ctx, params)
	if err != nil {
		return model.TransactionDetails{}, err
	}

	if err := s.accountService.Deposit(ctx, userID, amount, operationID(details)); err != nil {
		return model.TransactionDetails{}, fmt.Errorf("failed to update account balance: %w", err)
	}

	return details, nil
}

func (s *TransactionService) Withdraw(ctx context.Context, accountID uint64, amount int64) (model.TransactionDetails, error) {
	params := model.WithdrawParams{AccountID: accountID, Amount: amount}
	details, err := s.repo.Withdraw(ctx, params)
	if err != nil {
		return model.TransactionDetails{}, err
	}

	if err := s.accountService.Withdraw(ctx, accountID, amount, operationID(details)); err != nil {
		return model.TransactionDetails{}, fmt.Errorf("failed to update account balance: %w", err)
	}

	return details, nil
}

func (s *TransactionService) Transfer(ctx context.Context, userID, recipient uint64, amount int64) (model.TransactionDetails, error) {
	params := model.TransferParams{UserID: userID, Recipient: recipient, Amount: amount}
	details, err := s.repo.Transfer(ctx, params)
	if err != nil {
		return model.TransactionDetails{}, err
	}

	if err := s.accountService.Transfer(ctx, userID, recipient, amount, operationID(details)); err != nil {
		return model.TransactionDetails{}, fmt.Errorf("failed to update account balance: %w", err)
	}

	return details, nil
}

func operationID(details model.TransactionDetails) string {
	return fmt.Sprintf("transaction-%d", details.Transaction.ID)
}
