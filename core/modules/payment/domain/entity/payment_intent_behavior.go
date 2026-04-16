package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"

	sharedevents "go-socket/core/shared/contracts/events"
	"go-socket/core/shared/pkg/stackErr"
)

var (
	ErrPaymentProviderRequired          = errors.New("provider is required")
	ErrPaymentTransactionIDRequired     = errors.New("transaction_id is required")
	ErrPaymentAmountInvalid             = errors.New("amount must be greater than 0")
	ErrPaymentCurrencyRequired          = errors.New("currency is required")
	ErrPaymentClearingAccountKeyMissing = errors.New("clearing_account_key is required")
	ErrPaymentCreditAccountRequired     = errors.New("credit_account_id is required")
	ErrPaymentStatusInvalid             = errors.New("status is invalid")
	ErrPaymentProviderAmountMismatch    = errors.New("provider amount does not match reserved payment")
	ErrPaymentProviderCurrencyMismatch  = errors.New("provider currency does not match reserved payment")
	ErrPaymentProcessedProviderRequired = errors.New("provider is required")
	ErrPaymentProcessedKeyRequired      = errors.New("idempotency_key is required")
	ErrPaymentProcessedTxnRequired      = errors.New("transaction_id is required")
)

func NewProviderTopUpIntent(
	transactionID,
	provider string,
	amount int64,
	currency,
	beneficiaryAccountID string,
	now time.Time,
) (*PaymentIntent, error) {
	return newPaymentIntent(
		transactionID,
		provider,
		amount,
		currency,
		providerClearingAccountKey(provider),
		beneficiaryAccountID,
		now,
	)
}

func newPaymentIntent(transactionID, provider string, amount int64, currency, clearingAccountKey, creditAccountID string, now time.Time) (*PaymentIntent, error) {
	transactionID = strings.TrimSpace(transactionID)
	provider = strings.ToLower(strings.TrimSpace(provider))
	currency = strings.ToUpper(strings.TrimSpace(currency))
	clearingAccountKey = effectivePaymentClearingAccountKey(provider, clearingAccountKey)
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
	case clearingAccountKey == "":
		return nil, ErrPaymentClearingAccountKeyMissing
	case creditAccountID == "":
		return nil, ErrPaymentCreditAccountRequired
	}

	now = normalizePaymentTime(now)
	return &PaymentIntent{
		TransactionID:      transactionID,
		Provider:           provider,
		Amount:             amount,
		Currency:           currency,
		ClearingAccountKey: clearingAccountKey,
		CreditAccountID:    creditAccountID,
		Status:             PaymentStatusCreating,
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func providerClearingAccountKey(provider string) string {
	provider = strings.ToLower(strings.TrimSpace(provider))
	if provider == "" {
		return ""
	}
	return fmt.Sprintf("provider:%s", provider)
}

func effectivePaymentClearingAccountKey(provider, clearingAccountKey string) string {
	if clearingAccountKey = strings.TrimSpace(clearingAccountKey); clearingAccountKey != "" {
		return clearingAccountKey
	}
	return providerClearingAccountKey(provider)
}

func NormalizePaymentStatus(status string) string {
	normalized := strings.ToUpper(strings.TrimSpace(status))
	return ValidPaymentStatuses[normalized]
}

func NormalizePaymentStatusOrPending(status string) string {
	if normalized := NormalizePaymentStatus(status); normalized != "" {
		return normalized
	}
	return PaymentStatusPending
}

func (p *PaymentIntent) SetProviderState(externalRef, status string, updatedAt time.Time) error {
	_, err := p.transitionProviderState(externalRef, status, updatedAt)
	return err
}

func (p *PaymentIntent) TransitionProviderResult(result PaymentProviderResult, updatedAt time.Time) (PaymentTransition, error) {
	if p != nil {
		p.ensureWorkflowDefaults()
	}
	nextStatus := NormalizePaymentStatusOrPending(result.Status)
	if err := p.ValidateProviderResultForStatus(nextStatus, result.Amount, result.Currency); err != nil {
		return PaymentTransition{}, stackErr.Error(err)
	}

	return p.transitionProviderState(result.ExternalRef, nextStatus, updatedAt)
}

func (p *PaymentIntent) ApplyProviderResult(result PaymentProviderResult, updatedAt time.Time) error {
	_, err := p.TransitionProviderResult(result, updatedAt)
	return err
}

func (p *PaymentIntent) transitionProviderState(externalRef, status string, updatedAt time.Time) (PaymentTransition, error) {
	if p == nil {
		return PaymentTransition{}, ErrPaymentTransactionIDRequired
	}
	p.ensureWorkflowDefaults()

	normalizedStatus := NormalizePaymentStatus(status)
	if normalizedStatus == "" {
		return PaymentTransition{}, ErrPaymentStatusInvalid
	}

	transition := resolvePaymentTransition(NormalizePaymentStatusOrPending(p.Status), normalizedStatus)
	if transition.StateChanged {
		p.Status = transition.CurrentStatus
	}

	if externalRef = strings.TrimSpace(externalRef); externalRef != "" && externalRef != p.ExternalRef {
		p.ExternalRef = externalRef
		transition.ExternalRefChanged = true
	}

	if transition.StateChanged || transition.ExternalRefChanged {
		p.UpdatedAt = normalizePaymentTime(updatedAt)
	}

	return transition, nil
}

func resolvePaymentTransition(previousStatus, nextStatus string) PaymentTransition {
	transition := PaymentTransition{
		PreviousStatus: NormalizePaymentStatusOrPending(previousStatus),
		CurrentStatus:  NormalizePaymentStatusOrPending(previousStatus),
		Type:           PaymentTransitionNone,
	}

	if transition.CurrentStatus == nextStatus {
		return transition
	}

	switch transition.CurrentStatus {
	case PaymentStatusCreating:
		switch nextStatus {
		case PaymentStatusPending:
			transition.CurrentStatus = nextStatus
			transition.Type = PaymentTransitionPending
			transition.StateChanged = true
		case PaymentStatusSuccess:
			transition.CurrentStatus = nextStatus
			transition.Type = PaymentTransitionSucceeded
			transition.StateChanged = true
		case PaymentStatusFailed:
			transition.CurrentStatus = nextStatus
			transition.Type = PaymentTransitionFailed
			transition.StateChanged = true
		case PaymentStatusCancelled:
			transition.CurrentStatus = nextStatus
			transition.Type = PaymentTransitionCancelled
			transition.StateChanged = true
		default:
			transition.Ignored = true
		}
	case PaymentStatusPending:
		switch nextStatus {
		case PaymentStatusSuccess:
			transition.CurrentStatus = nextStatus
			transition.Type = PaymentTransitionSucceeded
			transition.StateChanged = true
		case PaymentStatusFailed:
			transition.CurrentStatus = nextStatus
			transition.Type = PaymentTransitionFailed
			transition.StateChanged = true
		case PaymentStatusCancelled:
			transition.CurrentStatus = nextStatus
			transition.Type = PaymentTransitionCancelled
			transition.StateChanged = true
		default:
			transition.Ignored = true
		}
	case PaymentStatusSuccess:
		switch nextStatus {
		case PaymentStatusRefunded:
			transition.CurrentStatus = nextStatus
			transition.Type = PaymentTransitionRefunded
			transition.StateChanged = true
		case PaymentStatusChargeback:
			transition.CurrentStatus = nextStatus
			transition.Type = PaymentTransitionChargeback
			transition.StateChanged = true
		default:
			transition.Ignored = true
		}
	case PaymentStatusFailed, PaymentStatusCancelled, PaymentStatusRefunded, PaymentStatusChargeback:
		transition.Ignored = true
	default:
		transition.Ignored = true
	}

	return transition
}

func (p *PaymentIntent) MarkCreateFailed(updatedAt time.Time) (PaymentTransition, error) {
	return p.transitionProviderState("", PaymentStatusFailed, updatedAt)
}

func (p *PaymentIntent) IsSucceeded() bool {
	return p != nil && p.Status == PaymentStatusSuccess
}

func (p *PaymentIntent) IsFailed() bool {
	return p != nil && p.Status == PaymentStatusFailed
}

func (p *PaymentIntent) IsCancelled() bool {
	return p != nil && p.Status == PaymentStatusCancelled
}

func (p *PaymentIntent) IsRefunded() bool {
	return p != nil && p.Status == PaymentStatusRefunded
}

func (p *PaymentIntent) IsChargeback() bool {
	return p != nil && p.Status == PaymentStatusChargeback
}

func (p *PaymentIntent) IsTerminal() bool {
	if p == nil {
		return false
	}

	switch p.Status {
	case PaymentStatusSuccess, PaymentStatusFailed, PaymentStatusCancelled, PaymentStatusRefunded, PaymentStatusChargeback:
		return true
	default:
		return false
	}
}

func (p *PaymentIntent) IsFinalized() bool {
	return p != nil && (p.IsSucceeded() || p.IsRefunded() || p.IsChargeback())
}

func (p *PaymentIntent) ShouldEmitCheckoutSessionCreated(checkoutURL string) bool {
	if p == nil {
		return false
	}
	return strings.TrimSpace(checkoutURL) != "" || strings.TrimSpace(p.ExternalRef) != ""
}

func (p *PaymentIntent) ValidateProviderResult(amount int64, currency string) error {
	return p.ValidateProviderResultForStatus(p.Status, amount, currency)
}

func (p *PaymentIntent) ValidateProviderResultForStatus(status string, amount int64, currency string) error {
	if p == nil {
		return ErrPaymentTransactionIDRequired
	}

	normalizedStatus := NormalizePaymentStatus(status)
	if normalizedStatus == "" {
		normalizedStatus = NormalizePaymentStatusOrPending(status)
	}

	switch normalizedStatus {
	case PaymentStatusRefunded, PaymentStatusChargeback:
		if amount != 0 && amount > p.Amount {
			return ErrPaymentProviderAmountMismatch
		}
	default:
		if amount != 0 && amount != p.Amount {
			return ErrPaymentProviderAmountMismatch
		}
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

func (p *PaymentIntent) TransitionIdempotencyKey(eventName string) string {
	return fmt.Sprintf("%s:%s", strings.TrimSpace(eventName), strings.TrimSpace(p.TransactionID))
}

func (p *PaymentIntent) BuildCreatedEventData(metadata map[string]string, createdAt time.Time) sharedevents.PaymentCreatedEvent {
	p.ensureWorkflowDefaults()
	occurredAt := normalizePaymentTime(createdAt)
	return sharedevents.PaymentCreatedEvent{
		PaymentID:          p.TransactionID,
		TransactionID:      p.TransactionID,
		Provider:           p.Provider,
		ClearingAccountKey: p.ClearingAccountKey,
		Amount:             p.Amount,
		Currency:           p.Currency,
		CreditAccountID:    p.CreditAccountID,
		Status:             p.Status,
		Metadata:           metadata,
		CreatedAt:          occurredAt,
	}
}

func (p *PaymentIntent) BuildCheckoutSessionCreatedEventData(checkoutURL string, occurredAt time.Time) sharedevents.PaymentCheckoutSessionCreatedEvent {
	p.ensureWorkflowDefaults()
	eventTime := normalizePaymentTime(occurredAt)
	return sharedevents.PaymentCheckoutSessionCreatedEvent{
		PaymentID:          p.TransactionID,
		TransactionID:      p.TransactionID,
		Provider:           p.Provider,
		ProviderPaymentRef: p.ExternalRef,
		CheckoutURL:        strings.TrimSpace(checkoutURL),
		Amount:             p.Amount,
		Currency:           p.Currency,
		Status:             p.Status,
		OccurredAt:         eventTime,
	}
}

func (p *PaymentIntent) BuildSucceededEventData(result PaymentProviderResult, occurredAt time.Time) sharedevents.PaymentSucceededEvent {
	p.ensureWorkflowDefaults()
	eventTime := normalizePaymentTime(occurredAt)
	return sharedevents.PaymentSucceededEvent{
		PaymentID:          p.TransactionID,
		TransactionID:      p.TransactionID,
		Provider:           p.Provider,
		ClearingAccountKey: p.ClearingAccountKey,
		ProviderEventID:    strings.TrimSpace(result.EventID),
		ProviderEventType:  strings.TrimSpace(result.EventType),
		ProviderPaymentRef: coalescePaymentValue(result.ExternalRef, p.ExternalRef),
		Amount:             p.Amount,
		Currency:           p.Currency,
		CreditAccountID:    p.CreditAccountID,
		IdempotencyKey:     p.TransitionIdempotencyKey(sharedevents.EventPaymentSucceeded),
		SucceededAt:        eventTime,
	}
}

func (p *PaymentIntent) BuildFailedEventData(result PaymentProviderResult, occurredAt time.Time) sharedevents.PaymentFailedEvent {
	p.ensureWorkflowDefaults()
	eventTime := normalizePaymentTime(occurredAt)
	return sharedevents.PaymentFailedEvent{
		PaymentID:          p.TransactionID,
		TransactionID:      p.TransactionID,
		Provider:           p.Provider,
		ProviderEventID:    strings.TrimSpace(result.EventID),
		ProviderEventType:  strings.TrimSpace(result.EventType),
		ProviderPaymentRef: coalescePaymentValue(result.ExternalRef, p.ExternalRef),
		Amount:             p.Amount,
		Currency:           p.Currency,
		Status:             NormalizePaymentStatusOrPending(result.Status),
		OccurredAt:         eventTime,
	}
}

func (p *PaymentIntent) BuildRefundedEventData(result PaymentProviderResult, occurredAt time.Time) sharedevents.PaymentRefundedEvent {
	p.ensureWorkflowDefaults()
	eventTime := normalizePaymentTime(occurredAt)
	current := p.CurrentProviderResult(result)
	return sharedevents.PaymentRefundedEvent{
		PaymentID:          p.TransactionID,
		TransactionID:      p.TransactionID,
		Provider:           p.Provider,
		ClearingAccountKey: p.ClearingAccountKey,
		ProviderEventID:    strings.TrimSpace(current.EventID),
		ProviderEventType:  strings.TrimSpace(current.EventType),
		ProviderPaymentRef: coalescePaymentValue(current.ExternalRef, p.ExternalRef),
		Amount:             paymentResultAmountOrDefault(current.Amount, p.Amount),
		Currency:           p.Currency,
		CreditAccountID:    p.CreditAccountID,
		IdempotencyKey:     p.TransitionIdempotencyKey(sharedevents.EventPaymentRefunded),
		RefundedAt:         eventTime,
	}
}

func (p *PaymentIntent) BuildChargebackEventData(result PaymentProviderResult, occurredAt time.Time) sharedevents.PaymentChargebackEvent {
	p.ensureWorkflowDefaults()
	eventTime := normalizePaymentTime(occurredAt)
	current := p.CurrentProviderResult(result)
	return sharedevents.PaymentChargebackEvent{
		PaymentID:          p.TransactionID,
		TransactionID:      p.TransactionID,
		Provider:           p.Provider,
		ClearingAccountKey: p.ClearingAccountKey,
		ProviderEventID:    strings.TrimSpace(current.EventID),
		ProviderEventType:  strings.TrimSpace(current.EventType),
		ProviderPaymentRef: coalescePaymentValue(current.ExternalRef, p.ExternalRef),
		Amount:             paymentResultAmountOrDefault(current.Amount, p.Amount),
		Currency:           p.Currency,
		CreditAccountID:    p.CreditAccountID,
		IdempotencyKey:     p.TransitionIdempotencyKey(sharedevents.EventPaymentChargeback),
		ChargedBackAt:      eventTime,
	}
}

func (p *PaymentIntent) NewProcessedEvent(result PaymentProviderResult, createdAt time.Time) (*ProcessedPaymentEvent, error) {
	return NewProcessedPaymentEvent(
		p.Provider,
		p.PaymentIdempotencyKey(result.EventID, result.ExternalRef),
		p.TransactionID,
		createdAt,
	)
}

func (p *PaymentIntent) NewProcessedTransitionEvent(eventName string, createdAt time.Time) (*ProcessedPaymentEvent, error) {
	return NewProcessedPaymentEvent(
		p.Provider,
		p.TransitionIdempotencyKey(eventName),
		p.TransactionID,
		createdAt,
	)
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

func coalescePaymentValue(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func paymentResultAmountOrDefault(amount int64, fallback int64) int64 {
	if amount != 0 {
		return amount
	}
	return fallback
}

func (p *PaymentIntent) ensureWorkflowDefaults() {
	if p == nil {
		return
	}
	p.TransactionID = strings.TrimSpace(p.TransactionID)
	p.Provider = strings.ToLower(strings.TrimSpace(p.Provider))
	p.ExternalRef = strings.TrimSpace(p.ExternalRef)
	p.Currency = strings.ToUpper(strings.TrimSpace(p.Currency))
	p.ClearingAccountKey = effectivePaymentClearingAccountKey(p.Provider, p.ClearingAccountKey)
	p.CreditAccountID = strings.TrimSpace(p.CreditAccountID)
	if p.Status = NormalizePaymentStatus(p.Status); p.Status == "" {
		p.Status = PaymentStatusCreating
	}
}

func (p *PaymentIntent) CurrentProviderResult(source PaymentProviderResult) PaymentProviderResult {
	if p == nil {
		return PaymentProviderResult{}
	}

	amount := source.Amount
	if amount == 0 {
		amount = p.Amount
	}

	return PaymentProviderResult{
		TransactionID: coalescePaymentValue(source.TransactionID, p.TransactionID),
		EventID:       strings.TrimSpace(source.EventID),
		EventType:     strings.TrimSpace(source.EventType),
		Status:        NormalizePaymentStatusOrPending(coalescePaymentValue(source.Status, p.Status)),
		Amount:        amount,
		Currency:      coalescePaymentValue(source.Currency, p.Currency),
		ExternalRef:   coalescePaymentValue(source.ExternalRef, p.ExternalRef),
	}
}
