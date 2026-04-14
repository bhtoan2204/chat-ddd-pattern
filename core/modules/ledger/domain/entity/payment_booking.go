package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"
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

type PaymentSucceededBookingInput struct {
	PaymentID       string
	TransactionID   string
	DebitAccountID  string
	CreditAccountID string
	Amount          int64
	IdempotencyKey  string
}

func NewPaymentSucceededBooking(input PaymentSucceededBookingInput) (*PaymentSucceededBooking, error) {
	paymentID := strings.TrimSpace(input.PaymentID)
	if paymentID == "" {
		paymentID = strings.TrimSpace(input.TransactionID)
	}
	if paymentID == "" {
		return nil, ErrPaymentBookingIDRequired
	}

	debitAccountID := strings.TrimSpace(input.DebitAccountID)
	creditAccountID := strings.TrimSpace(input.CreditAccountID)
	if debitAccountID == "" || creditAccountID == "" {
		return nil, ErrPaymentBookingAccountsRequired
	}
	if input.Amount <= 0 {
		return nil, ErrPaymentBookingAmountInvalid
	}

	idempotencyKey := strings.TrimSpace(input.IdempotencyKey)
	if idempotencyKey == "" {
		idempotencyKey = fmt.Sprintf("payment.succeeded:%s", paymentID)
	}

	return &PaymentSucceededBooking{
		PaymentID:       paymentID,
		DebitAccountID:  debitAccountID,
		CreditAccountID: creditAccountID,
		Amount:          input.Amount,
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
