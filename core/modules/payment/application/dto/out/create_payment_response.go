// CODE_GENERATOR - do not edit: response
package out

type CreatePaymentResponse struct {
	Provider      string `json:"provider,omitempty"`
	TransactionID string `json:"transaction_id,omitempty"`
	ExternalRef   string `json:"external_ref,omitempty"`
	Status        string `json:"status,omitempty"`
	CheckoutURL   string `json:"checkout_url,omitempty"`
}
