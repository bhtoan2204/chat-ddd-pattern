// CODE_GENERATOR - do not edit: response
package out

type CreateWithdrawalResponse struct {
	Provider       string `json:"provider,omitempty"`
	Workflow       string `json:"workflow,omitempty"`
	TransactionID  string `json:"transaction_id,omitempty"`
	ExternalRef    string `json:"external_ref,omitempty"`
	Amount         int64  `json:"amount,omitempty"`
	FeeAmount      int64  `json:"fee_amount,omitempty"`
	ProviderAmount int64  `json:"provider_amount,omitempty"`
	Status         string `json:"status,omitempty"`
}
