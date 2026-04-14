package aggregate

import (
	"errors"
	"fmt"
	"time"

	"go-socket/core/modules/ledger/domain/entity"
	"go-socket/core/shared/pkg/stackErr"
)

var ErrLedgerTransactionAggregateRequired = errors.New("ledger transaction aggregate is required")

type LedgerTransactionAggregate struct {
	transaction *entity.LedgerTransaction
}

func NewLedgerTransactionAggregate(transactionID string, entries []entity.LedgerEntryInput, now time.Time) (*LedgerTransactionAggregate, error) {
	transaction, err := entity.NewLedgerTransaction(transactionID, entries, now)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return &LedgerTransactionAggregate{transaction: transaction}, nil
}

func (a *LedgerTransactionAggregate) Snapshot() (*entity.LedgerTransaction, error) {
	if a == nil || a.transaction == nil {
		return nil, stackErr.Error(ErrLedgerTransactionAggregateRequired)
	}

	entries := make([]*entity.LedgerEntry, 0, len(a.transaction.Entries))
	for _, entry := range a.transaction.Entries {
		if entry == nil {
			continue
		}

		entryCopy := *entry
		entries = append(entries, &entryCopy)
	}

	return &entity.LedgerTransaction{
		TransactionID: a.transaction.TransactionID,
		CreatedAt:     a.transaction.CreatedAt,
		Entries:       entries,
	}, nil
}

func (a *LedgerTransactionAggregate) AssignEntryIDs(entryIDs []int64) error {
	if a == nil || a.transaction == nil {
		return stackErr.Error(ErrLedgerTransactionAggregateRequired)
	}
	if len(entryIDs) != len(a.transaction.Entries) {
		return stackErr.Error(fmt.Errorf("entry ids count mismatch"))
	}

	for idx, entryID := range entryIDs {
		if a.transaction.Entries[idx] == nil {
			return stackErr.Error(fmt.Errorf("entry at index %d is nil", idx))
		}
		a.transaction.Entries[idx].ID = entryID
	}

	return nil
}
