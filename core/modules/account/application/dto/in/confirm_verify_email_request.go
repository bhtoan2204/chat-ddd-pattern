// CODE_GENERATOR: request

package in

import "errors"

type ConfirmVerifyEmailRequest struct {
	Token string `json:"token" form:"token" binding:"required"`
}

func (r *ConfirmVerifyEmailRequest) Validate() error {
	if r.Token == "" {
		return errors.New("token is required")
	}
	return nil
}
