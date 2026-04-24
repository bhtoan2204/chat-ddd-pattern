// CODE_GENERATOR - do not edit: request

package in

import (
	"strings"
)

type ProcessStripeWebhookRequest struct {
	Signature string `json:"signature" form:"signature"`
	Payload   string `json:"payload" form:"payload"`
}

func (r *ProcessStripeWebhookRequest) Normalize() {
	r.Signature = strings.TrimSpace(r.Signature)
}

func (r *ProcessStripeWebhookRequest) Validate() error {
	r.Normalize()
	return nil
}
