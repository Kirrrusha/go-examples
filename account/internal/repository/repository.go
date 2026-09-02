package repository

import (
	"context"
	"errors"
	"fmt"

	"account/internal/model"
	"account/internal/repository/mapper"
	repomodel "account/internal/repository/model"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewRepository(db *gorm.DB, logger *zerolog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	userRepo := mapper.UserToRepoUser(user)

	res := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&userRepo)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to save user")
		return model.User{}, fmt.Errorf("failed to save user: %w", res.Error)
	}

	return mapper.RepoUserToUser(userRepo), nil
}

func (r *Repository) GetUser(ctx context.Context, userID uint64) (model.User, error) {
	var user repomodel.User

	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ? AND is_deleted = ?", userID, false).
		First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.User{}, fmt.Errorf("user not found")
	} else if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to get user by id")
		return model.User{}, res.Error
	}

	return mapper.RepoUserToUser(user), nil
}

func (r *Repository) GetUsers(ctx context.Context, limit int, offset int) ([]model.User, error) {
	var users []repomodel.User

	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("is_deleted = ?", false).
		Offset(offset).
		Limit(limit).
		Find(&users)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("users not found")
	} else if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to get users")
		return nil, res.Error
	}

	return mapper.RepoUsersToUsers(users), nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID uint64) error {
	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ?", userID).
		Update("is_deleted", true)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to delete user")
		return fmt.Errorf("failed to delete user: %w", res.Error)
	}

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, userID uint64, user model.UpdateUser) error {
	userRepo := mapper.UpdateUserToRepoUser(user)

	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ? AND is_deleted = ?", userID, false).
		Updates(userRepo)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to update user")
		return fmt.Errorf("failed to update user: %w", res.Error)
	}

	return nil
}

func (r *Repository) Deposit(ctx context.Context, userID uint64, amount int64) (int64, error) {
	return r.updateBalance(ctx, userID, amount)
}

func (r *Repository) Withdraw(ctx context.Context, userID uint64, amount int64) (int64, error) {
	return r.updateBalance(ctx, userID, -amount)
}

func (r *Repository) Transfer(ctx context.Context, userID uint64, recipientID uint64, amount int64) (int64, int64, error) {
	var userBalance int64
	var recipientBalance int64

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var sender repomodel.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND is_deleted = ?", userID, false).
			First(&sender).Error; err != nil {
			return fmt.Errorf("sender not found: %w", err)
		}
		if sender.Balance < amount {
			return fmt.Errorf("insufficient funds")
		}

		var recipient repomodel.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND is_deleted = ?", recipientID, false).
			First(&recipient).Error; err != nil {
			return fmt.Errorf("recipient not found: %w", err)
		}

		sender.Balance -= amount
		recipient.Balance += amount

		if err := tx.Model(&repomodel.User{}).Where("id = ?", sender.ID).Update("balance", sender.Balance).Error; err != nil {
			return err
		}
		if err := tx.Model(&repomodel.User{}).Where("id = ?", recipient.ID).Update("balance", recipient.Balance).Error; err != nil {
			return err
		}

		userBalance = sender.Balance
		recipientBalance = recipient.Balance
		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	return userBalance, recipientBalance, nil
}

func (r *Repository) GetBalance(ctx context.Context, userID uint64) (int64, error) {
	var user repomodel.User
	res := r.db.WithContext(ctx).
		Model(&repomodel.User{}).
		Where("id = ? AND is_deleted = ?", userID, false).
		First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0, fmt.Errorf("user not found")
	} else if res.Error != nil {
		return 0, res.Error
	}

	return user.Balance, nil
}

func (r *Repository) updateBalance(ctx context.Context, userID uint64, delta int64) (int64, error) {
	var balance int64

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user repomodel.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND is_deleted = ?", userID, false).
			First(&user).Error; err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		nextBalance := user.Balance + delta
		if nextBalance < 0 {
			return fmt.Errorf("insufficient funds")
		}

		if err := tx.Model(&repomodel.User{}).Where("id = ?", userID).Update("balance", nextBalance).Error; err != nil {
			return err
		}

		balance = nextBalance
		return nil
	})
	if err != nil {
		return 0, err
	}

	return balance, nil
}
