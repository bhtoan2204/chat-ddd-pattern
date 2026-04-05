// CODE_GENERATOR: request

package in

import "errors"

type RegisterRequest struct {
	DisplayName string `json:"display_name" form:"display_name" binding:"required"`
	Email       string `json:"email" form:"email" binding:"required,email"`
	Password    string `json:"password" form:"password" binding:"required"`
}

func (r *RegisterRequest) Validate() error {
	if r.DisplayName == "" {
		return errors.New("display_name is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
