// CODE_GENERATOR: request

package in

import "errors"

type UpdateRoomRequest struct {
	Id   string `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

func (r *UpdateRoomRequest) Validate() error {
	if r.Id == "" {
		return errors.New("id is required")
	}
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
