package repository

import (
	"testing"
	"time"

	"go-socket/core/modules/payment/domain/aggregate"
	"go-socket/core/modules/payment/infra/persistent/model"
)

func TestShouldCreatePaymentSnapshot(t *testing.T) {
	testCases := []struct {
		name        string
		baseVersion int
		newVersion  int
		expected    bool
	}{
		{name: "ignore zero version", baseVersion: 0, newVersion: 0, expected: false},
		{name: "create first snapshot", baseVersion: 0, newVersion: 1, expected: true},
		{name: "skip inside same bucket", baseVersion: 1, newVersion: 2, expected: false},
		{name: "create on interval boundary", baseVersion: 49, newVersion: 50, expected: true},
		{name: "create when interval crossed in same save", baseVersion: 49, newVersion: 51, expected: true},
		{name: "skip after boundary already captured", baseVersion: 50, newVersion: 51, expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := shouldCreatePaymentSnapshot(tc.baseVersion, tc.newVersion); actual != tc.expected {
				t.Fatalf("unexpected snapshot decision: got %v want %v", actual, tc.expected)
			}
		})
	}
}

func TestRestoreSnapshot(t *testing.T) {
	repo := &paymentBalanceAggregateRepoImpl{serializer: newPaymentBalanceSerializer()}

	original, err := aggregate.NewPaymentBalanceAggregate("account-1")
	if err != nil {
		t.Fatalf("new aggregate failed: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)
	if err := original.Deposit("tx-1", 100, now); err != nil {
		t.Fatalf("deposit failed: %v", err)
	}

	state, err := repo.serializer.Marshal(original)
	if err != nil {
		t.Fatalf("marshal snapshot failed: %v", err)
	}

	restored, err := aggregate.NewPaymentBalanceAggregate("account-1")
	if err != nil {
		t.Fatalf("new aggregate failed: %v", err)
	}

	err = repo.restoreSnapshot(restored, model.PaymentBalanceSnapshotModel{
		AggregateID: "account-1",
		Version:     1,
		State:       string(state),
	})
	if err != nil {
		t.Fatalf("restore snapshot failed: %v", err)
	}

	if restored.AccountID != "account-1" {
		t.Fatalf("unexpected account id: got %s", restored.AccountID)
	}
	if restored.Balance != 100 {
		t.Fatalf("unexpected balance: got %d want %d", restored.Balance, 100)
	}
	if !restored.CreatedAt.Equal(now) {
		t.Fatalf("unexpected created at: got %v want %v", restored.CreatedAt, now)
	}
	if restored.Root().Version() != 1 || restored.Root().BaseVersion() != 1 {
		t.Fatalf("unexpected root versions: got version=%d base=%d", restored.Root().Version(), restored.Root().BaseVersion())
	}
}
