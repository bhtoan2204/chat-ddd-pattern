package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"

	sharedevents "go-socket/core/shared/contracts/events"
)

var (
	ErrPaymentSucceededEventRequired   = errors.New("payment event is required")
	ErrPaymentBookingIDRequired        = errors.New("payment_id is required")
	ErrPaymentBookingAccountsRequired  = errors.New("debit_account_id and credit_account_id are required")
	ErrPaymentBookingAmountInvalid     = errors.New("amount must be positive")
	ErrPaymentBookingProviderRequired  = errors.New("provider is required")
	ErrPaymentBookingKeyRequired       = errors.New("idempotency_key is required")
	ErrPaymentBookingTransactionNeeded = errors.New("transaction_id is required")
)

type PaymentSucceededBooking struct {
	PaymentID       string
	DebitAccountID  string
	CreditAccountID string
	Amount          int64
	IdempotencyKey  string
}

func NewPaymentSucceededBooking(evt *sharedevents.PaymentSucceededEvent) (*PaymentSucceededBooking, error) {
	if evt == nil {
		return nil, ErrPaymentSucceededEventRequired
	}

	paymentID := strings.TrimSpace(evt.PaymentID)
	if paymentID == "" {
		paymentID = strings.TrimSpace(evt.TransactionID)
	}
	if paymentID == "" {
		return nil, ErrPaymentBookingIDRequired
	}

	debitAccountID := strings.TrimSpace(evt.DebitAccountID)
	creditAccountID := strings.TrimSpace(evt.CreditAccountID)
	if debitAccountID == "" || creditAccountID == "" {
		return nil, ErrPaymentBookingAccountsRequired
	}
	if evt.Amount <= 0 {
		return nil, ErrPaymentBookingAmountInvalid
	}

	idempotencyKey := strings.TrimSpace(evt.IdempotencyKey)
	if idempotencyKey == "" {
		idempotencyKey = fmt.Sprintf("payment.succeeded:%s", paymentID)
	}

	return &PaymentSucceededBooking{
		PaymentID:       paymentID,
		DebitAccountID:  debitAccountID,
		CreditAccountID: creditAccountID,
		Amount:          evt.Amount,
		IdempotencyKey:  idempotencyKey,
	}, nil
}

func (b *PaymentSucceededBooking) LedgerTransactionID() string {
	return fmt.Sprintf("payment:%s:succeeded", strings.TrimSpace(b.PaymentID))
}

func (b *PaymentSucceededBooking) LedgerEntries() []LedgerEntryInput {
	return []LedgerEntryInput{
		{AccountID: b.DebitAccountID, Amount: -b.Amount},
		{AccountID: b.CreditAccountID, Amount: b.Amount},
	}
}

func (b *PaymentSucceededBooking) ProcessedEvent(provider string, createdAt time.Time) (*ProcessedPaymentEvent, error) {
	provider = strings.TrimSpace(provider)
	if provider == "" {
		return nil, ErrPaymentBookingProviderRequired
	}
	if strings.TrimSpace(b.IdempotencyKey) == "" {
		return nil, ErrPaymentBookingKeyRequired
	}
	if strings.TrimSpace(b.PaymentID) == "" {
		return nil, ErrPaymentBookingTransactionNeeded
	}

	return &ProcessedPaymentEvent{
		Provider:       provider,
		IdempotencyKey: b.IdempotencyKey,
		TransactionID:  b.PaymentID,
		CreatedAt:      normalizeLedgerTime(createdAt),
	}, nil
}
