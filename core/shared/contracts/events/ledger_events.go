package events

import "time"

const (
	EventLedgerAccountTransferredToAccount = "EventLedgerAccountTransferredToAccount"
	EventLedgerPaymentReconciliationFailed = "EventLedgerPaymentReconciliationFailed"
)

type LedgerAccountTransferredToAccountEvent struct {
	TransactionID string    `json:"transaction_id"`
	ToAccountID   string    `json:"to_account_id"`
	Currency      string    `json:"currency"`
	Amount        int64     `json:"amount"`
	BookedAt      time.Time `json:"booked_at"`
}

type LedgerPaymentReconciliationFailedEvent struct {
	PaymentID          string    `json:"payment_id"`
	TransactionID      string    `json:"transaction_id"`
	Provider           string    `json:"provider"`
	ClearingAccountKey string    `json:"clearing_account_key"`
	CreditAccountID    string    `json:"credit_account_id,omitempty"`
	Currency           string    `json:"currency"`
	Amount             int64     `json:"amount"`
	FeeAmount          int64     `json:"fee_amount"`
	ProviderAmount     int64     `json:"provider_amount"`
	Reason             string    `json:"reason"`
	FailedAt           time.Time `json:"failed_at"`
}

type LedgerEntry struct {
	ID            int64     `json:"id"`
	TransactionID string    `json:"transaction_id"`
	AccountID     string    `json:"account_id"`
	Currency      string    `json:"currency"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type LedgerTransaction struct {
	TransactionID string
	Currency      string
	CreatedAt     time.Time
	Entries       []*LedgerEntry
}
