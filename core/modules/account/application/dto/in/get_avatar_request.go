// CODE_GENERATOR: request

package in

import "errors"

type GetAvatarRequest struct {
	AccountID string `json:"account_id" form:"account_id" binding:"required"`
}

func (r *GetAvatarRequest) Validate() error {
	if r.AccountID == "" {
		return errors.New("account_id is required")
	}
	return nil
}
