package entity

import "time"

const (
	PaymentAggregateType    = "payment"
	PaymentStatusCreating   = "CREATING"
	PaymentStatusPending    = "PENDING"
	PaymentStatusSuccess    = "SUCCESS"
	PaymentStatusFailed     = "FAILED"
	PaymentStatusCancelled  = "CANCELLED"
	PaymentStatusRefunded   = "REFUNDED"
	PaymentStatusChargeback = "CHARGEBACK"
)

type PaymentTransitionType string

const (
	PaymentTransitionNone       PaymentTransitionType = "none"
	PaymentTransitionPending    PaymentTransitionType = "pending"
	PaymentTransitionSucceeded  PaymentTransitionType = "succeeded"
	PaymentTransitionFailed     PaymentTransitionType = "failed"
	PaymentTransitionCancelled  PaymentTransitionType = "cancelled"
	PaymentTransitionRefunded   PaymentTransitionType = "refunded"
	PaymentTransitionChargeback PaymentTransitionType = "chargeback"
)

type PaymentTransition struct {
	PreviousStatus     string
	CurrentStatus      string
	Type               PaymentTransitionType
	StateChanged       bool
	ExternalRefChanged bool
	Ignored            bool
}

type PaymentIntent struct {
	TransactionID      string
	Provider           string
	ExternalRef        string
	Amount             int64
	Currency           string
	ClearingAccountKey string
	CreditAccountID    string
	Status             string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type PaymentProviderResult struct {
	TransactionID string
	EventID       string
	EventType     string
	Status        string
	Amount        int64
	Currency      string
	ExternalRef   string
}

type ProcessedPaymentEvent struct {
	Provider       string
	IdempotencyKey string
	TransactionID  string
	CreatedAt      time.Time
}
