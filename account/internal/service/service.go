package service

import (
	"context"
	"time"

	"account/internal/model"
	"github.com/rs/zerolog"
)

type AccountService struct {
	repo   Repository
	logger *zerolog.Logger
}

func New(repo Repository, logger *zerolog.Logger) *AccountService {
	return &AccountService{repo: repo, logger: logger}
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
