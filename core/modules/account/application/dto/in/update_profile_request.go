// CODE_GENERATOR: request

package in

import "errors"

type UpdateProfileRequest struct {
	DisplayName     string  `json:"display_name" form:"display_name" binding:"required"`
	Username        *string `json:"username,omitempty" form:"username"`
	AvatarObjectKey *string `json:"avatar_object_key,omitempty" form:"avatar_object_key"`
}

func (r *UpdateProfileRequest) Validate() error {
	if r.DisplayName == "" {
		return errors.New("display_name is required")
	}
	return nil
}
