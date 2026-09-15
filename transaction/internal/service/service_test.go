package service

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/rs/zerolog"
	"transaction/internal/model"
)

type transactionRepositoryStub struct {
	getTransactionsFn       func(context.Context, model.GetTransactionsParams) ([]model.Transaction, error)
	getTransactionDetailsFn func(context.Context, uint64) (model.TransactionDetails, error)
	depositFn               func(context.Context, model.DepositParams) (model.TransactionDetails, error)
	withdrawFn              func(context.Context, model.WithdrawParams) (model.TransactionDetails, error)
	transferFn              func(context.Context, model.TransferParams) (model.TransactionDetails, error)
	updateStatusFn          func(context.Context, uint64, model.TransactionStatus) error
}

func (s *transactionRepositoryStub) GetTransactions(ctx context.Context, params model.GetTransactionsParams) ([]model.Transaction, error) {
	return s.getTransactionsFn(ctx, params)
}
func (s *transactionRepositoryStub) GetTransactionDetails(ctx context.Context, id uint64) (model.TransactionDetails, error) {
	return s.getTransactionDetailsFn(ctx, id)
}
func (s *transactionRepositoryStub) Deposit(ctx context.Context, params model.DepositParams) (model.TransactionDetails, error) {
	return s.depositFn(ctx, params)
}
func (s *transactionRepositoryStub) Withdraw(ctx context.Context, params model.WithdrawParams) (model.TransactionDetails, error) {
	return s.withdrawFn(ctx, params)
}
func (s *transactionRepositoryStub) Transfer(ctx context.Context, params model.TransferParams) (model.TransactionDetails, error) {
	return s.transferFn(ctx, params)
}
func (s *transactionRepositoryStub) UpdateTransactionStatus(ctx context.Context, id uint64, status model.TransactionStatus) error {
	return s.updateStatusFn(ctx, id, status)
}

type transactionPublisherStub struct {
	topic string
	key   string
	data  interface{}
	err   error
	calls int
}

func (s *transactionPublisherStub) Publish(_ context.Context, topic, key string, data interface{}) error {
	s.topic, s.key, s.data = topic, key, data
	s.calls++
	return s.err
}

type accountServiceStub struct{}

func (accountServiceStub) Deposit(context.Context, uint64, int64, string) error  { return nil }
func (accountServiceStub) Withdraw(context.Context, uint64, int64, string) error { return nil }
func (accountServiceStub) Transfer(context.Context, uint64, uint64, int64, string) error {
	return nil
}

func newTransactionServiceForTest(repo Repository, publisher KafkaPublisher) *TransactionService {
	logger := zerolog.New(io.Discard)
	return New(repo, accountServiceStub{}, publisher, &logger)
}

func TestTransactionOperationsPublishRequests(t *testing.T) {
	tests := []struct {
		name      string
		operation func(*TransactionService) (model.TransactionDetails, error)
		prepare   func(*testing.T, *transactionRepositoryStub)
		wantType  string
		wantUser  uint64
		wantTo    uint64
		wantID    uint64
	}{
		{
			name: "deposit",
			operation: func(service *TransactionService) (model.TransactionDetails, error) {
				return service.Deposit(context.Background(), 7, 1250)
			},
			prepare: func(t *testing.T, repo *transactionRepositoryStub) {
				repo.depositFn = func(_ context.Context, params model.DepositParams) (model.TransactionDetails, error) {
					if params.UserID != 7 || params.Amount != 1250 {
						t.Fatalf("unexpected deposit params: %+v", params)
					}
					return transactionDetails(21, 7, 1250, model.TransactionTypeDeposit), nil
				}
			},
			wantType: "deposit", wantUser: 7, wantID: 21,
		},
		{
			name: "withdraw",
			operation: func(service *TransactionService) (model.TransactionDetails, error) {
				return service.Withdraw(context.Background(), 7, 500)
			},
			prepare: func(t *testing.T, repo *transactionRepositoryStub) {
				repo.withdrawFn = func(_ context.Context, params model.WithdrawParams) (model.TransactionDetails, error) {
					if params.AccountID != 7 || params.Amount != 500 {
						t.Fatalf("unexpected withdraw params: %+v", params)
					}
					return transactionDetails(22, 7, 500, model.TransactionTypeWithdraw), nil
				}
			},
			wantType: "withdraw", wantUser: 7, wantID: 22,
		},
		{
			name: "transfer",
			operation: func(service *TransactionService) (model.TransactionDetails, error) {
				return service.Transfer(context.Background(), 7, 8, 300)
			},
			prepare: func(t *testing.T, repo *transactionRepositoryStub) {
				repo.transferFn = func(_ context.Context, params model.TransferParams) (model.TransactionDetails, error) {
					if params.UserID != 7 || params.Recipient != 8 || params.Amount != 300 {
						t.Fatalf("unexpected transfer params: %+v", params)
					}
					return transactionDetails(23, 7, 300, model.TransactionTypeTransfer), nil
				}
			},
			wantType: "transfer", wantUser: 7, wantTo: 8, wantID: 23,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &transactionRepositoryStub{}
			test.prepare(t, repo)
			publisher := &transactionPublisherStub{}
			service := newTransactionServiceForTest(repo, publisher)

			result, err := test.operation(service)
			if err != nil {
				t.Fatalf("operation returned error: %v", err)
			}
			if result.Transaction.ID != test.wantID {
				t.Fatalf("unexpected transaction: %+v", result.Transaction)
			}
			if publisher.calls != 1 || publisher.topic != "transaction_data" || publisher.key != "7" {
				t.Fatalf("unexpected publish: calls=%d topic=%q key=%q", publisher.calls, publisher.topic, publisher.key)
			}
			request, ok := publisher.data.(map[string]interface{})
			if !ok {
				t.Fatalf("published value has type %T", publisher.data)
			}
			if request["request_type"] != test.wantType || request["user_id"] != test.wantUser || request["operation_id"] != test.wantID {
				t.Fatalf("unexpected request: %#v", request)
			}
			if test.wantTo != 0 && request["recipient_id"] != test.wantTo {
				t.Fatalf("unexpected recipient: %#v", request["recipient_id"])
			}
		})
	}
}

func TestTransactionOperationMarksFailedWhenPublishFails(t *testing.T) {
	wantErr := errors.New("kafka unavailable")
	var updatedID uint64
	var updatedStatus model.TransactionStatus
	repo := &transactionRepositoryStub{
		depositFn: func(context.Context, model.DepositParams) (model.TransactionDetails, error) {
			return transactionDetails(31, 7, 100, model.TransactionTypeDeposit), nil
		},
		updateStatusFn: func(_ context.Context, id uint64, status model.TransactionStatus) error {
			updatedID, updatedStatus = id, status
			return nil
		},
	}
	service := newTransactionServiceForTest(repo, &transactionPublisherStub{err: wantErr})

	_, err := service.Deposit(context.Background(), 7, 100)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected publish error, got %v", err)
	}
	if updatedID != 31 || updatedStatus != model.TransactionStatusFailed {
		t.Fatalf("unexpected status update: id=%d status=%s", updatedID, updatedStatus)
	}
}

func TestHandleAccountResponseUpdatesStatus(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus model.TransactionStatus
	}{
		{name: "completed", body: `{"request_type":"deposit","user_id":7,"operation_id":41,"result":true}`, wantStatus: model.TransactionStatusCompleted},
		{name: "failed", body: `{"request_type":"withdraw","user_id":7,"operation_id":41,"result":false}`, wantStatus: model.TransactionStatusFailed},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotStatus model.TransactionStatus
			repo := &transactionRepositoryStub{
				updateStatusFn: func(_ context.Context, id uint64, status model.TransactionStatus) error {
					if id != 41 {
						t.Fatalf("unexpected transaction id: %d", id)
					}
					gotStatus = status
					return nil
				},
			}
			service := newTransactionServiceForTest(repo, &transactionPublisherStub{})
			if err := service.HandleAccountResponse(context.Background(), "transaction_response", "7", []byte(test.body)); err != nil {
				t.Fatalf("HandleAccountResponse returned error: %v", err)
			}
			if gotStatus != test.wantStatus {
				t.Fatalf("got status %s, want %s", gotStatus, test.wantStatus)
			}
		})
	}
}

func TestHandleAccountResponseRejectsInvalidMessages(t *testing.T) {
	service := newTransactionServiceForTest(&transactionRepositoryStub{}, &transactionPublisherStub{})
	for _, body := range []string{`{`, `{"request_type":"deposit"}`} {
		if err := service.HandleAccountResponse(context.Background(), "transaction_response", "", []byte(body)); err == nil {
			t.Fatalf("expected error for %q", body)
		}
	}
}

func TestGetTransactionsWithDetailsSkipsUnavailableDetails(t *testing.T) {
	repo := &transactionRepositoryStub{
		getTransactionsFn: func(context.Context, model.GetTransactionsParams) ([]model.Transaction, error) {
			return []model.Transaction{{ID: 51}, {ID: 52}}, nil
		},
		getTransactionDetailsFn: func(_ context.Context, id uint64) (model.TransactionDetails, error) {
			if id == 51 {
				return transactionDetails(51, 7, 100, model.TransactionTypeDeposit), nil
			}
			return model.TransactionDetails{}, errors.New("not found")
		},
	}
	service := newTransactionServiceForTest(repo, &transactionPublisherStub{})
	details, err := service.GetTransactionsWithDetails(context.Background(), model.GetTransactionsParams{})
	if err != nil {
		t.Fatalf("GetTransactionsWithDetails returned error: %v", err)
	}
	if len(details) != 1 || details[0].Transaction.ID != 51 {
		t.Fatalf("unexpected details: %+v", details)
	}
}

func transactionDetails(id, userID uint64, amount int64, transactionType model.TransactionType) model.TransactionDetails {
	return model.TransactionDetails{Transaction: model.Transaction{
		ID: id, UserID: userID, Amount: amount, Status: model.TransactionStatusPending, Type: transactionType,
	}}
}
