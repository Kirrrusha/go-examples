package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"transaction/internal/model"
	"transaction/internal/repository/mapper"
	repomodel "transaction/internal/repository/model"
)

type Repository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewRepository(db *gorm.DB, logger *zerolog.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

func (r *Repository) GetTransactions(ctx context.Context, params model.GetTransactionsParams) ([]model.Transaction, error) {
	var transactions []repomodel.Transaction
	query := r.db.WithContext(ctx).Model(&repomodel.Transaction{})

	if params.UserID != nil {
		query = query.Where("user_id = ?", *params.UserID)
	}
	if params.Type != nil {
		query = query.Where("type = ?", *params.Type)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.DateFrom != nil {
		query = query.Where("created_at >= ?", *params.DateFrom)
	}
	if params.DateTo != nil {
		query = query.Where("created_at <= ?", *params.DateTo)
	}

	res := query.Offset(params.Offset).Limit(params.Limit).Order("created_at DESC").Find(&transactions)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to get transactions")
		return nil, fmt.Errorf("failed to get transactions: %w", res.Error)
	}

	return mapper.RepoTransactionsToTransactions(transactions), nil
}

func (r *Repository) GetTransactionDetails(ctx context.Context, transactionID uint64) (model.TransactionDetails, error) {
	var transaction repomodel.Transaction
	res := r.db.WithContext(ctx).
		Model(&repomodel.Transaction{}).
		Where("id = ?", transactionID).
		First(&transaction)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return model.TransactionDetails{}, fmt.Errorf("transaction not found")
	} else if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to get transaction")
		return model.TransactionDetails{}, fmt.Errorf("failed to get transaction: %w", res.Error)
	}

	var entries []repomodel.TransactionEntry
	res = r.db.WithContext(ctx).
		Model(&repomodel.TransactionEntry{}).
		Where("transaction_id = ?", transactionID).
		Order("created_at ASC").
		Find(&entries)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("failed to get transaction entries")
		return model.TransactionDetails{}, fmt.Errorf("failed to get transaction entries: %w", res.Error)
	}

	return model.TransactionDetails{
		Transaction: mapper.RepoTransactionToTransaction(transaction),
		Entries:     mapper.RepoTransactionEntriesToTransactionEntries(entries),
	}, nil
}

func (r *Repository) Deposit(ctx context.Context, params model.DepositParams) (model.TransactionDetails, error) {
	return r.createOperation(ctx, operationParams{
		userID: params.UserID,
		amount: params.Amount,
		txType: model.TransactionTypeDeposit,
		entries: []model.TransactionEntry{
			{
				AccountID: params.UserID,
				Direction: model.TransactionEntryDirectionDebit,
				Amount:    params.Amount,
			},
		},
	})
}

func (r *Repository) Withdraw(ctx context.Context, params model.WithdrawParams) (model.TransactionDetails, error) {
	return r.createOperation(ctx, operationParams{
		userID: params.AccountID,
		amount: params.Amount,
		txType: model.TransactionTypeWithdraw,
		entries: []model.TransactionEntry{
			{
				AccountID: params.AccountID,
				Direction: model.TransactionEntryDirectionCredit,
				Amount:    params.Amount,
			},
		},
	})
}

func (r *Repository) Transfer(ctx context.Context, params model.TransferParams) (model.TransactionDetails, error) {
	return r.createOperation(ctx, operationParams{
		userID: params.UserID,
		amount: params.Amount,
		txType: model.TransactionTypeTransfer,
		entries: []model.TransactionEntry{
			{
				AccountID: params.UserID,
				Direction: model.TransactionEntryDirectionCredit,
				Amount:    params.Amount,
			},
			{
				AccountID: params.Recipient,
				Direction: model.TransactionEntryDirectionDebit,
				Amount:    params.Amount,
			},
		},
	})
}

type operationParams struct {
	userID  uint64
	amount  int64
	txType  model.TransactionType
	entries []model.TransactionEntry
}

func (r *Repository) createOperation(ctx context.Context, params operationParams) (model.TransactionDetails, error) {
	var result model.TransactionDetails

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		transaction := mapper.TransactionToRepoTransaction(model.Transaction{
			UserID:    params.userID,
			Amount:    params.amount,
			Status:    model.TransactionStatusPending,
			Type:      params.txType,
			CreatedAt: now,
			UpdatedAt: now,
		})

		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&transaction).Error; err != nil {
			r.logger.Err(err).Msg("failed to create transaction")
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		entries := make([]model.TransactionEntry, len(params.entries))
		for i, entry := range params.entries {
			entry.TransactionID = transaction.ID
			entry.CreatedAt = now
			entry.UpdatedAt = now

			entryRepo := mapper.TransactionEntryToRepoTransactionEntry(entry)
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&entryRepo).Error; err != nil {
				r.logger.Err(err).Msg("failed to create transaction entry")
				return fmt.Errorf("failed to create transaction entry: %w", err)
			}

			entries[i] = mapper.RepoTransactionEntryToTransactionEntry(entryRepo)
		}

		updateData := mapper.UpdateTransactionToRepoTransaction(model.UpdateTransaction{
			Status: model.TransactionStatusCompleted,
		})
		if err := tx.Model(&repomodel.Transaction{}).Where("id = ?", transaction.ID).Updates(&updateData).Error; err != nil {
			r.logger.Err(err).Msg("failed to update transaction status")
			return fmt.Errorf("failed to update transaction status: %w", err)
		}

		transaction.Status = string(model.TransactionStatusCompleted)
		result = model.TransactionDetails{
			Transaction: mapper.RepoTransactionToTransaction(transaction),
			Entries:     entries,
		}
		return nil
	})
	if err != nil {
		return model.TransactionDetails{}, err
	}

	return result, nil
}
