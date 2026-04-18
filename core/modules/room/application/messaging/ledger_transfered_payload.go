package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	sharedevents "wechat-clone/core/shared/contracts/events"
	"wechat-clone/core/shared/pkg/stackErr"
)

type ledgerTransferPayload struct {
	TransactionID     string
	SenderAccountID   string
	ReceiverAccountID string
	Currency          string
	AmountMinor       int64
	CreatedAt         time.Time
}

func decodeLedgerAccountTransferPayload(ctx context.Context, raw json.RawMessage) (*ledgerTransferPayload, error) {
	payloadAny, err := decodeEventPayload(ctx, sharedevents.EventLedgerAccountTransferredToAccount, raw)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	payload, ok := payloadAny.(*sharedevents.LedgerTransaction)
	if !ok {
		return nil, stackErr.Error(fmt.Errorf("invalid payload type for event %s", sharedevents.EventLedgerAccountTransferredToAccount))
	}
	if payload == nil {
		return nil, stackErr.Error(fmt.Errorf("ledger transfer payload is required"))
	}
	if strings.TrimSpace(payload.TransactionID) == "" {
		return nil, stackErr.Error(fmt.Errorf("ledger transfer transaction_id is required"))
	}
	if len(payload.Entries) != 2 {
		return nil, stackErr.Error(fmt.Errorf("ledger transfer payload must contain exactly 2 entries"))
	}

	debitEntry, creditEntry, err := splitTransferEntries(payload.Entries)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	senderID := strings.TrimSpace(debitEntry.AccountID)
	receiverID := strings.TrimSpace(creditEntry.AccountID)
	if senderID == "" || receiverID == "" {
		return nil, stackErr.Error(fmt.Errorf("ledger transfer payload account ids are required"))
	}

	currency, err := resolveTransferCurrency(payload, debitEntry, creditEntry)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	amountMinor := debitEntry.Amount * -1
	if amountMinor != creditEntry.Amount {
		return nil, stackErr.Error(fmt.Errorf("ledger transfer payload amounts are unbalanced"))
	}

	return &ledgerTransferPayload{
		TransactionID:     strings.TrimSpace(payload.TransactionID),
		SenderAccountID:   senderID,
		ReceiverAccountID: receiverID,
		Currency:          currency,
		AmountMinor:       amountMinor,
		CreatedAt:         payload.CreatedAt.UTC(),
	}, nil
}

func splitTransferEntries(entries []*sharedevents.LedgerEntry) (*sharedevents.LedgerEntry, *sharedevents.LedgerEntry, error) {
	var debitEntry *sharedevents.LedgerEntry
	var creditEntry *sharedevents.LedgerEntry

	for _, entry := range entries {
		if entry == nil {
			return nil, nil, stackErr.Error(fmt.Errorf("ledger transfer payload contains nil entry"))
		}

		switch {
		case entry.Amount < 0:
			if debitEntry != nil {
				return nil, nil, stackErr.Error(fmt.Errorf("ledger transfer payload must contain exactly one debit entry"))
			}
			debitEntry = entry
		case entry.Amount > 0:
			if creditEntry != nil {
				return nil, nil, stackErr.Error(fmt.Errorf("ledger transfer payload must contain exactly one credit entry"))
			}
			creditEntry = entry
		default:
			return nil, nil, stackErr.Error(fmt.Errorf("ledger transfer payload entry amount must be non-zero"))
		}
	}

	if debitEntry == nil || creditEntry == nil {
		return nil, nil, stackErr.Error(fmt.Errorf("ledger transfer payload must contain one debit and one credit entry"))
	}

	return debitEntry, creditEntry, nil
}

func resolveTransferCurrency(
	payload *sharedevents.LedgerTransaction,
	debitEntry *sharedevents.LedgerEntry,
	creditEntry *sharedevents.LedgerEntry,
) (string, error) {
	currency := strings.ToUpper(strings.TrimSpace(payload.Currency))
	if currency == "" {
		currency = strings.ToUpper(strings.TrimSpace(debitEntry.Currency))
	}
	if currency == "" {
		currency = strings.ToUpper(strings.TrimSpace(creditEntry.Currency))
	}
	if currency == "" {
		return "", stackErr.Error(fmt.Errorf("ledger transfer currency is required"))
	}

	if !strings.EqualFold(strings.TrimSpace(debitEntry.Currency), currency) || !strings.EqualFold(strings.TrimSpace(creditEntry.Currency), currency) {
		return "", stackErr.Error(fmt.Errorf("ledger transfer entries currency mismatch"))
	}

	return currency, nil
}

func transferMessageID(transactionID string) string {
	return "ledger-transfer:" + strings.TrimSpace(transactionID)
}

func formatTransferMessageBody(currency string, amountMinor int64) string {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	return currency + " " + strconv.FormatInt(amountMinor, 10)
}
