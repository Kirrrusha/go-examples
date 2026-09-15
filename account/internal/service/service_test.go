package service

import (
	"context"
	"errors"
	"io"
	"testing"

	"account/internal/model"
	"github.com/rs/zerolog"
)

type accountRepositoryStub struct {
	depositFn  func(context.Context, uint64, int64) (int64, error)
	withdrawFn func(context.Context, uint64, int64) (int64, error)
	transferFn func(context.Context, uint64, uint64, int64) (int64, int64, error)
}

func (s *accountRepositoryStub) CreateUser(context.Context, model.User) (model.User, error) {
	return model.User{}, nil
}
func (s *accountRepositoryStub) GetUser(context.Context, uint64) (model.User, error) {
	return model.User{}, nil
}
func (s *accountRepositoryStub) GetUsers(context.Context, int, int) ([]model.User, error) {
	return nil, nil
}
func (s *accountRepositoryStub) DeleteUser(context.Context, uint64) error { return nil }
func (s *accountRepositoryStub) UpdateUser(context.Context, uint64, model.UpdateUser) error {
	return nil
}
func (s *accountRepositoryStub) Deposit(ctx context.Context, userID uint64, amount int64) (int64, error) {
	return s.depositFn(ctx, userID, amount)
}
func (s *accountRepositoryStub) Withdraw(ctx context.Context, userID uint64, amount int64) (int64, error) {
	return s.withdrawFn(ctx, userID, amount)
}
func (s *accountRepositoryStub) Transfer(ctx context.Context, userID, recipientID uint64, amount int64) (int64, int64, error) {
	return s.transferFn(ctx, userID, recipientID, amount)
}
func (s *accountRepositoryStub) GetBalance(context.Context, uint64) (int64, error) {
	return 0, nil
}

type accountPublisherStub struct {
	topic string
	key   string
	data  interface{}
	err   error
	calls int
}

func (s *accountPublisherStub) Publish(_ context.Context, topic, key string, data interface{}) error {
	s.topic, s.key, s.data = topic, key, data
	s.calls++
	return s.err
}

func newAccountServiceForTest(repo Repository, publisher KafkaPublisher) *AccountService {
	logger := zerolog.New(io.Discard)
	return New(repo, publisher, &logger)
}

func TestHandleTransactionRoutesSuccessfulOperations(t *testing.T) {
	tests := []struct {
		name string
		body string
		want func(*testing.T, *accountRepositoryStub)
	}{
		{
			name: "deposit",
			body: `{"request_type":"deposit","user_id":7,"amount":1250,"operation_id":11}`,
			want: func(t *testing.T, repo *accountRepositoryStub) {
				repo.depositFn = func(_ context.Context, userID uint64, amount int64) (int64, error) {
					if userID != 7 || amount != 1250 {
						t.Fatalf("unexpected deposit arguments: user=%d amount=%d", userID, amount)
					}
					return 2250, nil
				}
			},
		},
		{
			name: "withdraw",
			body: `{"request_type":"withdraw","user_id":7,"amount":500,"operation_id":12}`,
			want: func(t *testing.T, repo *accountRepositoryStub) {
				repo.withdrawFn = func(_ context.Context, userID uint64, amount int64) (int64, error) {
					if userID != 7 || amount != 500 {
						t.Fatalf("unexpected withdraw arguments: user=%d amount=%d", userID, amount)
					}
					return 500, nil
				}
			},
		},
		{
			name: "transfer",
			body: `{"request_type":"transfer","user_id":7,"recipient_id":8,"amount":300,"operation_id":13}`,
			want: func(t *testing.T, repo *accountRepositoryStub) {
				repo.transferFn = func(_ context.Context, userID, recipientID uint64, amount int64) (int64, int64, error) {
					if userID != 7 || recipientID != 8 || amount != 300 {
						t.Fatalf("unexpected transfer arguments: from=%d to=%d amount=%d", userID, recipientID, amount)
					}
					return 700, 1300, nil
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &accountRepositoryStub{}
			test.want(t, repo)
			publisher := &accountPublisherStub{}
			service := newAccountServiceForTest(repo, publisher)

			if err := service.HandleTransaction(context.Background(), "transaction_data", "7", []byte(test.body)); err != nil {
				t.Fatalf("HandleTransaction returned error: %v", err)
			}
			if publisher.calls != 1 || publisher.topic != "transaction_response" || publisher.key != "7" {
				t.Fatalf("unexpected publish: calls=%d topic=%q key=%q", publisher.calls, publisher.topic, publisher.key)
			}
			response, ok := publisher.data.(TransactionResponse)
			if !ok {
				t.Fatalf("published value has type %T", publisher.data)
			}
			if !response.Result || response.UserID != 7 || response.OperationID == 0 || response.RequestType != test.name {
				t.Fatalf("unexpected response: %+v", response)
			}
		})
	}
}

func TestHandleTransactionPublishesFailedOperation(t *testing.T) {
	wantErr := errors.New("insufficient funds")
	repo := &accountRepositoryStub{
		withdrawFn: func(context.Context, uint64, int64) (int64, error) { return 0, wantErr },
	}
	publisher := &accountPublisherStub{}
	service := newAccountServiceForTest(repo, publisher)

	err := service.HandleTransaction(context.Background(), "transaction_data", "7", []byte(`{"request_type":"withdraw","user_id":7,"amount":500,"operation_id":14}`))
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected operation error, got %v", err)
	}
	response, ok := publisher.data.(TransactionResponse)
	if !ok || response.Result {
		t.Fatalf("expected failed response, got %#v", publisher.data)
	}
}

func TestHandleTransactionRejectsInvalidMessages(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "invalid json", body: `{`},
		{name: "missing fields", body: `{"request_type":"deposit"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			publisher := &accountPublisherStub{}
			service := newAccountServiceForTest(&accountRepositoryStub{}, publisher)
			if err := service.HandleTransaction(context.Background(), "transaction_data", "", []byte(test.body)); err == nil {
				t.Fatal("expected validation error")
			}
			if publisher.calls != 0 {
				t.Fatalf("unexpected publish count: %d", publisher.calls)
			}
		})
	}
}

func TestHandleTransactionUnknownTypeAndPublishError(t *testing.T) {
	t.Run("unknown type publishes negative response", func(t *testing.T) {
		publisher := &accountPublisherStub{}
		service := newAccountServiceForTest(&accountRepositoryStub{}, publisher)
		err := service.HandleTransaction(context.Background(), "transaction_data", "7", []byte(`{"request_type":"refund","user_id":7,"operation_id":15}`))
		if err == nil {
			t.Fatal("expected unknown request type error")
		}
		response, ok := publisher.data.(TransactionResponse)
		if !ok || response.Result {
			t.Fatalf("expected negative response, got %#v", publisher.data)
		}
	})

	t.Run("publish error is returned", func(t *testing.T) {
		wantErr := errors.New("kafka unavailable")
		repo := &accountRepositoryStub{
			depositFn: func(context.Context, uint64, int64) (int64, error) { return 100, nil },
		}
		publisher := &accountPublisherStub{err: wantErr}
		service := newAccountServiceForTest(repo, publisher)
		err := service.HandleTransaction(context.Background(), "transaction_data", "7", []byte(`{"request_type":"deposit","user_id":7,"amount":100,"operation_id":16}`))
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected publish error, got %v", err)
		}
	})
}
