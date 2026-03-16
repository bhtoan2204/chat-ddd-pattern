package aggregate

import (
	"errors"
	"testing"
	"time"
)

func TestPaymentBalanceAggregateDepositAndWithdraw(t *testing.T) {
	agg, err := NewPaymentBalanceAggregate("account-1")
	if err != nil {
		t.Fatalf("new aggregate failed: %v", err)
	}

	now := time.Now().UTC()
	if err := agg.Deposit("tx-deposit", 100, now); err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	agg.Update()

	if err := agg.Withdraw("tx-withdraw", 40, now.Add(time.Minute)); err != nil {
		t.Fatalf("withdraw failed: %v", err)
	}

	if agg.Balance != 60 {
		t.Fatalf("unexpected balance: got %d want %d", agg.Balance, 60)
	}
}

func TestPaymentBalanceAggregateWithdrawInsufficientBalance(t *testing.T) {
	agg, err := NewPaymentBalanceAggregate("account-1")
	if err != nil {
		t.Fatalf("new aggregate failed: %v", err)
	}

	err = agg.Withdraw("tx-withdraw", 10, time.Now().UTC())
	if !errors.Is(err, ErrInsufficientBalance) {
		t.Fatalf("unexpected error: %v", err)
	}
}
