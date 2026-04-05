// CODE_GENERATOR: request

package in

import "errors"

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" form:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" form:"new_password" binding:"required"`
}

func (r *ChangePasswordRequest) Validate() error {
	if r.CurrentPassword == "" {
		return errors.New("current_password is required")
	}
	if r.NewPassword == "" {
		return errors.New("new_password is required")
	}
	return nil
}
