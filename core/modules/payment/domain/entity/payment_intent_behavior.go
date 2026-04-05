package entity

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrPaymentProviderRequired          = errors.New("provider is required")
	ErrPaymentTransactionIDRequired     = errors.New("transaction_id is required")
	ErrPaymentAmountInvalid             = errors.New("amount must be greater than 0")
	ErrPaymentCurrencyRequired          = errors.New("currency is required")
	ErrPaymentDebitAccountRequired      = errors.New("debit_account_id is required")
	ErrPaymentCreditAccountRequired     = errors.New("credit_account_id is required")
	ErrPaymentAccountsMustDiffer        = errors.New("debit_account_id and credit_account_id must be different")
	ErrPaymentStatusInvalid             = errors.New("status is invalid")
	ErrPaymentProviderAmountMismatch    = errors.New("provider amount does not match reserved payment")
	ErrPaymentProviderCurrencyMismatch  = errors.New("provider currency does not match reserved payment")
	ErrPaymentProcessedProviderRequired = errors.New("provider is required")
	ErrPaymentProcessedKeyRequired      = errors.New("idempotency_key is required")
	ErrPaymentProcessedTxnRequired      = errors.New("transaction_id is required")
)

func NewPaymentIntent(transactionID, provider string, amount int64, currency, debitAccountID, creditAccountID string, now time.Time) (*PaymentIntent, error) {
	transactionID = strings.TrimSpace(transactionID)
	provider = strings.ToLower(strings.TrimSpace(provider))
	currency = strings.ToUpper(strings.TrimSpace(currency))
	debitAccountID = strings.TrimSpace(debitAccountID)
	creditAccountID = strings.TrimSpace(creditAccountID)

	switch {
	case provider == "":
		return nil, ErrPaymentProviderRequired
	case transactionID == "":
		return nil, ErrPaymentTransactionIDRequired
	case amount <= 0:
		return nil, ErrPaymentAmountInvalid
	case currency == "":
		return nil, ErrPaymentCurrencyRequired
	case debitAccountID == "":
		return nil, ErrPaymentDebitAccountRequired
	case creditAccountID == "":
		return nil, ErrPaymentCreditAccountRequired
	case debitAccountID == creditAccountID:
		return nil, ErrPaymentAccountsMustDiffer
	}

	now = normalizePaymentTime(now)
	return &PaymentIntent{
		TransactionID:   transactionID,
		Provider:        provider,
		Amount:          amount,
		Currency:        currency,
		DebitAccountID:  debitAccountID,
		CreditAccountID: creditAccountID,
		Status:          PaymentStatusCreating,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func NormalizePaymentStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case PaymentStatusCreating:
		return PaymentStatusCreating
	case PaymentStatusPending:
		return PaymentStatusPending
	case PaymentStatusSuccess:
		return PaymentStatusSuccess
	case PaymentStatusFailed:
		return PaymentStatusFailed
	default:
		return ""
	}
}

func NormalizePaymentStatusOrPending(status string) string {
	if normalized := NormalizePaymentStatus(status); normalized != "" {
		return normalized
	}
	return PaymentStatusPending
}

func (p *PaymentIntent) SetProviderState(externalRef, status string, updatedAt time.Time) error {
	if p == nil {
		return ErrPaymentTransactionIDRequired
	}

	normalizedStatus := NormalizePaymentStatus(status)
	if normalizedStatus == "" {
		return ErrPaymentStatusInvalid
	}

	if externalRef = strings.TrimSpace(externalRef); externalRef != "" {
		p.ExternalRef = externalRef
	}
	p.Status = normalizedStatus
	p.UpdatedAt = normalizePaymentTime(updatedAt)
	return nil
}

func (p *PaymentIntent) ValidateProviderResult(amount int64, currency string) error {
	if p == nil {
		return ErrPaymentTransactionIDRequired
	}
	if amount != 0 && amount != p.Amount {
		return ErrPaymentProviderAmountMismatch
	}
	if currency = strings.TrimSpace(currency); currency != "" && !strings.EqualFold(currency, p.Currency) {
		return ErrPaymentProviderCurrencyMismatch
	}
	return nil
}

func (p *PaymentIntent) PaymentIdempotencyKey(eventID, externalRef string) string {
	if eventID = strings.TrimSpace(eventID); eventID != "" {
		return eventID
	}
	if externalRef = strings.TrimSpace(externalRef); externalRef != "" {
		return externalRef
	}
	if externalRef = strings.TrimSpace(p.ExternalRef); externalRef != "" {
		return externalRef
	}
	return strings.TrimSpace(p.TransactionID)
}

func NewProcessedPaymentEvent(provider, idempotencyKey, transactionID string, createdAt time.Time) (*ProcessedPaymentEvent, error) {
	provider = strings.TrimSpace(provider)
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	transactionID = strings.TrimSpace(transactionID)

	switch {
	case provider == "":
		return nil, ErrPaymentProcessedProviderRequired
	case idempotencyKey == "":
		return nil, ErrPaymentProcessedKeyRequired
	case transactionID == "":
		return nil, ErrPaymentProcessedTxnRequired
	}

	return &ProcessedPaymentEvent{
		Provider:       provider,
		IdempotencyKey: idempotencyKey,
		TransactionID:  transactionID,
		CreatedAt:      normalizePaymentTime(createdAt),
	}, nil
}

func normalizePaymentTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now().UTC()
	}
	return value.UTC()
}
