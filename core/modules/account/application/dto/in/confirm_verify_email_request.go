// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type ConfirmVerifyEmailRequest struct {
	Token string `json:"token" form:"token" binding:"required"`
}

func (r *ConfirmVerifyEmailRequest) Normalize() {
	r.Token = strings.TrimSpace(r.Token)
}

func (r *ConfirmVerifyEmailRequest) Validate() error {
	r.Normalize()
	if r.Token == "" {
		return stackErr.Error(errors.New("token is required"))
	}
	return nil
}
