// CODE_GENERATOR - do not edit: response
package out

type CreateQuoteResponse struct {
	QuoteID      string `json:"quote_id,omitempty"`
	FromCurrency string `json:"from_currency,omitempty"`
	ToCurrency   int64  `json:"to_currency,omitempty"`
	FromAmount   int64  `json:"from_amount,omitempty"`
	ToAmount     string `json:"to_amount,omitempty"`
	CustomerRate string `json:"customer_rate,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
}
