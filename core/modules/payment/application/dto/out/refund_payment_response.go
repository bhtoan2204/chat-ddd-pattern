// CODE_GENERATOR - do not edit: response
package out

type RefundPaymentResponse struct {
	Provider      string                    `json:"provider,omitempty"`
	TransactionID string                    `json:"transaction_id,omitempty"`
	ExternalRef   string                    `json:"external_ref,omitempty"`
	Status        string                    `json:"status,omitempty"`
	Duplicate     bool                      `json:"duplicate,omitempty"`
	Events        []PaymentIntegrationEvent `json:"events,omitempty"`
}
