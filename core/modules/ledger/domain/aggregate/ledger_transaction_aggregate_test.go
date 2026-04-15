package aggregate

import (
	"testing"

	"go-socket/core/modules/ledger/domain/entity"
)

func TestLedgerTransactionAggregateCreateBuildsSnapshot(t *testing.T) {
	aggregate, err := NewLedgerTransactionAggregate(" ledger-tx-1 ")
	if err != nil {
		t.Fatalf("NewLedgerTransactionAggregate() error = %v", err)
	}
	if err := aggregate.Create([]entity.LedgerEntryInput{
		{AccountID: " debit ", Currency: "vnd", Amount: -250},
		{AccountID: " credit ", Currency: "VND", Amount: 250},
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	snapshot, err := aggregate.Snapshot()
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}

	if snapshot.TransactionID != "ledger-tx-1" {
		t.Fatalf("unexpected transaction id: %s", snapshot.TransactionID)
	}
	if snapshot.Currency != "VND" {
		t.Fatalf("unexpected transaction currency: %s", snapshot.Currency)
	}
	if len(snapshot.Entries) != 2 {
		t.Fatalf("unexpected entries count: %d", len(snapshot.Entries))
	}
	if snapshot.Entries[0].AccountID != "debit" || snapshot.Entries[1].AccountID != "credit" {
		t.Fatalf("unexpected account ids: %+v", snapshot.Entries)
	}
	if snapshot.Entries[0].Currency != "VND" || snapshot.Entries[1].Currency != "VND" {
		t.Fatalf("unexpected entry currencies: %+v", snapshot.Entries)
	}
	if aggregate.Root().Version() != 1 {
		t.Fatalf("unexpected aggregate version: %d", aggregate.Root().Version())
	}
}
