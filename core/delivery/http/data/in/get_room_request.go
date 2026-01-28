// CODE_GENERATOR: request

package in

type GetRoomRequest struct {
	Id string `json:"id" form:"id"`
}

func (r *GetRoomRequest) Validate() error {
	return nil
}
