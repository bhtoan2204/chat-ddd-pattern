// CODE_GENERATOR - do not edit: response
package out

type TransactionTransactionResponse struct {
	TransactionID string                `json:"transaction_id,omitempty"`
	Currency      string                `json:"currency,omitempty"`
	CreatedAt     string                `json:"created_at,omitempty"`
	Entries       []LedgerEntryResponse `json:"entries,omitempty"`
}
