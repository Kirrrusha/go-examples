package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"transaction/internal/model"
)

type TransactionService struct {
	repo           Repository
	accountService AccountService
	kafkaPublisher KafkaPublisher
	logger         *zerolog.Logger
}

func New(repo Repository, accountService AccountService, kafkaPublisher KafkaPublisher, logger *zerolog.Logger) *TransactionService {
	return &TransactionService{repo: repo, accountService: accountService, kafkaPublisher: kafkaPublisher, logger: logger}
}

type AccountResponse struct {
	RequestType string `json:"request_type"`
	UserID      uint64 `json:"user_id"`
	OperationID uint64 `json:"operation_id"`
	Result      bool   `json:"result"`
}

type Repository interface {
	GetTransactions(context.Context, model.GetTransactionsParams) ([]model.Transaction, error)
	GetTransactionDetails(context.Context, uint64) (model.TransactionDetails, error)
	Deposit(context.Context, model.DepositParams) (model.TransactionDetails, error)
	Withdraw(context.Context, model.WithdrawParams) (model.TransactionDetails, error)
	Transfer(context.Context, model.TransferParams) (model.TransactionDetails, error)
	UpdateTransactionStatus(context.Context, uint64, model.TransactionStatus) error
}

type KafkaPublisher interface {
	Publish(ctx context.Context, topic, key string, data interface{}) error
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

	request := map[string]interface{}{
		"request_type": "deposit", "user_id": userID, "amount": amount,
		"operation_id": details.Transaction.ID, "timestamp": time.Now().UTC(),
	}
	if err := s.kafkaPublisher.Publish(ctx, "transaction_data", fmt.Sprintf("%d", userID), request); err != nil {
		s.markFailed(ctx, details.Transaction.ID)
		return model.TransactionDetails{}, fmt.Errorf("failed to publish deposit request: %w", err)
	}

	return details, nil
}

func (s *TransactionService) Withdraw(ctx context.Context, accountID uint64, amount int64) (model.TransactionDetails, error) {
	params := model.WithdrawParams{AccountID: accountID, Amount: amount}
	details, err := s.repo.Withdraw(ctx, params)
	if err != nil {
		return model.TransactionDetails{}, err
	}

	request := map[string]interface{}{
		"request_type": "withdraw", "user_id": accountID, "amount": amount,
		"operation_id": details.Transaction.ID, "timestamp": time.Now().UTC(),
	}
	if err := s.kafkaPublisher.Publish(ctx, "transaction_data", fmt.Sprintf("%d", accountID), request); err != nil {
		s.markFailed(ctx, details.Transaction.ID)
		return model.TransactionDetails{}, fmt.Errorf("failed to publish withdraw request: %w", err)
	}

	return details, nil
}

func (s *TransactionService) Transfer(ctx context.Context, userID, recipient uint64, amount int64) (model.TransactionDetails, error) {
	params := model.TransferParams{UserID: userID, Recipient: recipient, Amount: amount}
	details, err := s.repo.Transfer(ctx, params)
	if err != nil {
		return model.TransactionDetails{}, err
	}

	request := map[string]interface{}{
		"request_type": "transfer", "user_id": userID, "recipient_id": recipient,
		"amount": amount, "operation_id": details.Transaction.ID, "timestamp": time.Now().UTC(),
	}
	if err := s.kafkaPublisher.Publish(ctx, "transaction_data", fmt.Sprintf("%d", userID), request); err != nil {
		s.markFailed(ctx, details.Transaction.ID)
		return model.TransactionDetails{}, fmt.Errorf("failed to publish transfer request: %w", err)
	}

	return details, nil
}

func (s *TransactionService) HandleAccountResponse(ctx context.Context, topic, key string, data []byte) error {
	var response AccountResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("failed to unmarshal account response: %w", err)
	}
	if response.RequestType == "" || response.UserID == 0 || response.OperationID == 0 {
		return fmt.Errorf("account response is missing required fields")
	}

	status := model.TransactionStatusFailed
	if response.Result {
		status = model.TransactionStatusCompleted
	}
	if err := s.repo.UpdateTransactionStatus(ctx, response.OperationID, status); err != nil {
		return fmt.Errorf("failed to update transaction status to %s: %w", status, err)
	}
	s.logger.Info().Str("request_type", response.RequestType).Uint64("user_id", response.UserID).
		Uint64("operation_id", response.OperationID).Str("status", string(status)).Msg("processed account response")
	return nil
}

func (s *TransactionService) markFailed(ctx context.Context, transactionID uint64) {
	if err := s.repo.UpdateTransactionStatus(ctx, transactionID, model.TransactionStatusFailed); err != nil {
		s.logger.Error().Err(err).Uint64("transaction_id", transactionID).Msg("failed to mark transaction as failed")
	}
}
