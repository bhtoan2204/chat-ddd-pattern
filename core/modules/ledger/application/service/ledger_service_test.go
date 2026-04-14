package service

import (
	"context"
	"errors"
	"testing"

	ledgeraggregate "go-socket/core/modules/ledger/domain/aggregate"
	"go-socket/core/modules/ledger/domain/entity"
	ledgerrepos "go-socket/core/modules/ledger/domain/repos"

	"go.uber.org/mock/gomock"
)

func TestLedgerServiceCreateTransaction(t *testing.T) {
	t.Run("valid transaction", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		baseRepo := ledgerrepos.NewMockRepos(ctrl)
		txRepos := ledgerrepos.NewMockRepos(ctrl)
		aggregateRepo := ledgerrepos.NewMockLedgerTransactionAggregateRepository(ctrl)

		baseRepo.EXPECT().
			WithTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ledgerrepos.Repos) error) error {
				return fn(txRepos)
			})
		txRepos.EXPECT().LedgerTransactionAggregateRepository().Return(aggregateRepo)
		aggregateRepo.EXPECT().
			Save(gomock.Any(), gomock.AssignableToTypeOf(&ledgeraggregate.LedgerTransactionAggregate{})).
			DoAndReturn(func(_ context.Context, aggregate *ledgeraggregate.LedgerTransactionAggregate) error {
				transaction, err := aggregate.Snapshot()
				if err != nil {
					t.Fatalf("Snapshot() error = %v", err)
				}
				if transaction.TransactionID != "ledger-tx-valid" {
					t.Fatalf("expected transaction id ledger-tx-valid, got %s", transaction.TransactionID)
				}
				if len(transaction.Entries) != 2 {
					t.Fatalf("expected 2 entries, got %d", len(transaction.Entries))
				}
				return aggregate.AssignEntryIDs([]int64{1, 2})
			})

		service := NewLedgerService(baseRepo)

		transaction, err := service.CreateTransaction(context.Background(), CreateTransactionCommand{
			TransactionID: "ledger-tx-valid",
			Entries: []CreateTransactionEntryCommand{
				{AccountID: "acc-a", Amount: -100},
				{AccountID: "acc-b", Amount: 100},
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if transaction.TransactionID != "ledger-tx-valid" {
			t.Fatalf("expected transaction id ledger-tx-valid, got %s", transaction.TransactionID)
		}
		if len(transaction.Entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(transaction.Entries))
		}
		if transaction.Entries[0].ID != 1 || transaction.Entries[1].ID != 2 {
			t.Fatalf("expected persisted entry ids [1 2], got [%d %d]", transaction.Entries[0].ID, transaction.Entries[1].ID)
		}
	})

	t.Run("invalid transaction", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		baseRepo := ledgerrepos.NewMockRepos(ctrl)
		service := NewLedgerService(baseRepo)

		_, err := service.CreateTransaction(context.Background(), CreateTransactionCommand{
			TransactionID: "ledger-tx-invalid",
			Entries: []CreateTransactionEntryCommand{
				{AccountID: "acc-a", Amount: -100},
				{AccountID: "acc-b", Amount: 50},
			},
		})
		if !errors.Is(err, ErrValidation) {
			t.Fatalf("expected validation error, got %v", err)
		}
	})

	t.Run("duplicate transaction", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		baseRepo := ledgerrepos.NewMockRepos(ctrl)
		txRepos := ledgerrepos.NewMockRepos(ctrl)
		aggregateRepo := ledgerrepos.NewMockLedgerTransactionAggregateRepository(ctrl)

		baseRepo.EXPECT().
			WithTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(ledgerrepos.Repos) error) error {
				return fn(txRepos)
			})
		txRepos.EXPECT().LedgerTransactionAggregateRepository().Return(aggregateRepo)
		aggregateRepo.EXPECT().
			Save(gomock.Any(), gomock.AssignableToTypeOf(&ledgeraggregate.LedgerTransactionAggregate{})).
			Return(ledgerrepos.ErrDuplicate)

		service := NewLedgerService(baseRepo)

		_, err := service.CreateTransaction(context.Background(), CreateTransactionCommand{
			TransactionID: "ledger-tx-dup",
			Entries: []CreateTransactionEntryCommand{
				{AccountID: "acc-a", Amount: -100},
				{AccountID: "acc-b", Amount: 100},
			},
		})
		if !errors.Is(err, ErrDuplicateTransaction) {
			t.Fatalf("expected duplicate transaction error, got %v", err)
		}
	})
}

func TestLedgerServiceRecordPaymentSucceeded(t *testing.T) {
	ctrl := gomock.NewController(t)
	baseRepo := ledgerrepos.NewMockRepos(ctrl)
	txRepos := ledgerrepos.NewMockRepos(ctrl)
	aggregateRepo := ledgerrepos.NewMockLedgerTransactionAggregateRepository(ctrl)
	paymentRepo := ledgerrepos.NewMockPaymentRepository(ctrl)

	baseRepo.EXPECT().
		WithTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(ledgerrepos.Repos) error) error {
			return fn(txRepos)
		})
	txRepos.EXPECT().PaymentRepository().Return(paymentRepo).Times(2)
	paymentRepo.EXPECT().IsProcessed(gomock.Any(), paymentSucceededSource, "payment.succeeded:pay-1").Return(false, nil)
	txRepos.EXPECT().LedgerTransactionAggregateRepository().Return(aggregateRepo)
	aggregateRepo.EXPECT().
		Save(gomock.Any(), gomock.AssignableToTypeOf(&ledgeraggregate.LedgerTransactionAggregate{})).
		DoAndReturn(func(_ context.Context, aggregate *ledgeraggregate.LedgerTransactionAggregate) error {
			transaction, err := aggregate.Snapshot()
			if err != nil {
				t.Fatalf("Snapshot() error = %v", err)
			}
			if transaction.TransactionID != "payment:pay-1:succeeded" {
				t.Fatalf("expected payment booking transaction id, got %s", transaction.TransactionID)
			}
			if len(transaction.Entries) != 2 {
				t.Fatalf("expected 2 ledger entries, got %d", len(transaction.Entries))
			}
			return aggregate.AssignEntryIDs([]int64{10, 11})
		})
	paymentRepo.EXPECT().
		MarkProcessed(gomock.Any(), gomock.AssignableToTypeOf(&entity.ProcessedPaymentEvent{})).
		DoAndReturn(func(_ context.Context, event *entity.ProcessedPaymentEvent) error {
			if event.Provider != paymentSucceededSource {
				t.Fatalf("expected provider %s, got %s", paymentSucceededSource, event.Provider)
			}
			if event.TransactionID != "pay-1" {
				t.Fatalf("expected processed payment id pay-1, got %s", event.TransactionID)
			}
			if event.IdempotencyKey != "payment.succeeded:pay-1" {
				t.Fatalf("expected idempotency key payment.succeeded:pay-1, got %s", event.IdempotencyKey)
			}
			return nil
		})

	service := NewLedgerService(baseRepo)

	err := service.RecordPaymentSucceeded(context.Background(), RecordPaymentSucceededCommand{
		PaymentID:       "pay-1",
		DebitAccountID:  "wallet:pending",
		CreditAccountID: "wallet:available",
		Amount:          100,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
