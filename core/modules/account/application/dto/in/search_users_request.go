// CODE_GENERATOR - do not edit: request

package in

import (
	"errors"
	"strings"
	"wechat-clone/core/shared/pkg/stackErr"
)

type SearchUsersRequest struct {
	Q      string `json:"q" form:"q" binding:"required"`
	Limit  int    `json:"limit" form:"limit"`
	Offset int    `json:"offset" form:"offset"`
}

func (r *SearchUsersRequest) Normalize() {
	r.Q = strings.TrimSpace(r.Q)
}

func (r *SearchUsersRequest) Validate() error {
	r.Normalize()
	if r.Q == "" {
		return stackErr.Error(errors.New("q is required"))
	}
	return nil
}
