// CODE_GENERATOR: request

package in

type DeleteRoomRequest struct {
	Id string `json:"id" form:"id"`
}

func (r *DeleteRoomRequest) Validate() error {
	return nil
}
