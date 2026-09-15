package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"account/internal/model"
	"github.com/rs/zerolog"
)

type AccountService struct {
	repo           Repository
	kafkaPublisher KafkaPublisher
	logger         *zerolog.Logger
}

func New(repo Repository, kafkaPublisher KafkaPublisher, logger *zerolog.Logger) *AccountService {
	return &AccountService{repo: repo, kafkaPublisher: kafkaPublisher, logger: logger}
}

type TransactionRequest struct {
	RequestType string `json:"request_type"`
	UserID      uint64 `json:"user_id"`
	Amount      int64  `json:"amount"`
	OperationID uint64 `json:"operation_id"`
	RecipientID uint64 `json:"recipient_id"`
}

type TransactionResponse struct {
	RequestType string `json:"request_type"`
	UserID      uint64 `json:"user_id"`
	OperationID uint64 `json:"operation_id"`
	Result      bool   `json:"result"`
}

type Repository interface {
	CreateUser(context.Context, model.User) (model.User, error)
	GetUser(context.Context, uint64) (model.User, error)
	GetUsers(context.Context, int, int) ([]model.User, error)
	DeleteUser(context.Context, uint64) error
	UpdateUser(context.Context, uint64, model.UpdateUser) error
	Deposit(context.Context, uint64, int64) (int64, error)
	Withdraw(context.Context, uint64, int64) (int64, error)
	Transfer(context.Context, uint64, uint64, int64) (int64, int64, error)
	GetBalance(context.Context, uint64) (int64, error)
}

type KafkaPublisher interface {
	Publish(ctx context.Context, topic, key string, data interface{}) error
}

func (s *AccountService) CreateUser(ctx context.Context, newUser model.CreateUser) (model.User, error) {
	now := time.Now()
	user := model.User{
		Login:      newUser.Login,
		Email:      newUser.Email,
		Phone:      newUser.Phone,
		FirstName:  newUser.FirstName,
		LastName:   newUser.LastName,
		MiddleName: newUser.MiddleName,
		Age:        newUser.Age,
		Balance:    newUser.Balance,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *AccountService) GetUsers(ctx context.Context, limit int, offset int) ([]model.User, error) {
	return s.repo.GetUsers(ctx, limit, offset)
}

func (s *AccountService) GetUser(ctx context.Context, userID uint64) (model.User, error) {
	return s.repo.GetUser(ctx, userID)
}

func (s *AccountService) DeleteUser(ctx context.Context, userID uint64) error {
	return s.repo.DeleteUser(ctx, userID)
}

func (s *AccountService) UpdateUser(ctx context.Context, userID uint64, user model.UpdateUser) error {
	return s.repo.UpdateUser(ctx, userID, user)
}

func (s *AccountService) Deposit(ctx context.Context, userID uint64, amount int64) (int64, error) {
	return s.repo.Deposit(ctx, userID, amount)
}

func (s *AccountService) Withdraw(ctx context.Context, userID uint64, amount int64) (int64, error) {
	return s.repo.Withdraw(ctx, userID, amount)
}

func (s *AccountService) Transfer(ctx context.Context, userID uint64, recipientID uint64, amount int64) (int64, int64, error) {
	return s.repo.Transfer(ctx, userID, recipientID, amount)
}

func (s *AccountService) GetBalance(ctx context.Context, userID uint64) (int64, error) {
	return s.repo.GetBalance(ctx, userID)
}

func (s *AccountService) HandleTransaction(ctx context.Context, topic, key string, data []byte) error {
	var request TransactionRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return fmt.Errorf("failed to unmarshal transaction request: %w", err)
	}
	if request.RequestType == "" || request.UserID == 0 || request.OperationID == 0 {
		return fmt.Errorf("transaction request is missing required fields")
	}

	var err error
	switch request.RequestType {
	case "deposit":
		_, err = s.Deposit(ctx, request.UserID, request.Amount)
	case "withdraw":
		_, err = s.Withdraw(ctx, request.UserID, request.Amount)
	case "transfer":
		_, _, err = s.Transfer(ctx, request.UserID, request.RecipientID, request.Amount)
	default:
		err = fmt.Errorf("unknown request type: %s", request.RequestType)
	}

	response := TransactionResponse{
		RequestType: request.RequestType,
		UserID:      request.UserID,
		OperationID: request.OperationID,
		Result:      err == nil,
	}
	if publishErr := s.kafkaPublisher.Publish(ctx, "transaction_response", strconv.FormatUint(request.UserID, 10), response); publishErr != nil {
		return fmt.Errorf("failed to publish transaction response: %w", publishErr)
	}
	if err != nil {
		s.logger.Error().Err(err).Uint64("user_id", request.UserID).Uint64("operation_id", request.OperationID).Msg("transaction request failed")
		return fmt.Errorf("failed to handle transaction request: %w", err)
	}
	return nil
}
