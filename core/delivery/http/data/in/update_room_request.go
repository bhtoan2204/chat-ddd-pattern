// CODE_GENERATOR: request

package in

type UpdateRoomRequest struct {
	Id   string `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

func (r *UpdateRoomRequest) Validate() error {
	return nil
}
