// CODE_GENERATOR: request

package in

type CreateRoomRequest struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	RoomType    string `json:"room_type" form:"room_type"`
}

func (r *CreateRoomRequest) Validate() error {
	return nil
}
