// CODE_GENERATOR: request

package in

import (
	"errors"
	"go-socket/core/shared/pkg/stackErr"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

func (r *LoginRequest) Normalize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (r *LoginRequest) Validate() error {
	r.Normalize()
	if r.Email == "" {
		return stackErr.Error(errors.New("email is required"))
	}
	if r.Password == "" {
		return stackErr.Error(errors.New("password is required"))
	}
	return nil
}
