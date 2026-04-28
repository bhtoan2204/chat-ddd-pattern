package aggregate

import (
	"time"

	sharedevents "wechat-clone/core/shared/contracts/events"
)

type PaymentCreatedEvent = sharedevents.PaymentCreatedEvent

type PaymentWithdrawalRequestedEvent = sharedevents.PaymentWithdrawalRequestedEvent

type PaymentCheckoutSessionCreatedEvent = sharedevents.PaymentCheckoutSessionCreatedEvent

type PaymentSucceededEvent = sharedevents.PaymentSucceededEvent

type PaymentFailedEvent = sharedevents.PaymentFailedEvent

type PaymentRefundedEvent = sharedevents.PaymentRefundedEvent

type PaymentChargebackEvent = sharedevents.PaymentChargebackEvent

type PaymentProviderStateChangedEvent struct {
	TransactionID      string    `json:"transaction_id"`
	Provider           string    `json:"provider"`
	ProviderPaymentRef string    `json:"provider_payment_ref"`
	PreviousStatus     string    `json:"previous_status"`
	Status             string    `json:"status"`
	ProviderEventID    string    `json:"provider_event_id"`
	ProviderEventType  string    `json:"provider_event_type"`
	Amount             int64     `json:"amount"`
	Currency           string    `json:"currency"`
	OccurredAt         time.Time `json:"occurred_at"`
	StateChanged       bool      `json:"state_changed"`
	ExternalRefChanged bool      `json:"external_ref_changed"`
}
