// CODE_GENERATOR: response
package out

type ListRoomsResponse struct {
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Rooms []RoomResponse `json:"rooms"`
}

type RoomResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RoomType    string `json:"room_type"`
	OwnerId     string `json:"owner_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
