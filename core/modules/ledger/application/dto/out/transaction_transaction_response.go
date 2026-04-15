// CODE_GENERATOR - do not edit: response
package out

type TransactionTransactionResponse struct {
	TransactionID string                `json:"transaction_id"`
	Currency      string                `json:"currency"`
	CreatedAt     string                `json:"created_at"`
	Entries       []LedgerEntryResponse `json:"entries"`
}
