// CODE_GENERATOR: response
package out

type CreateMessageResponse struct {
	Id        string `json:"id"`
	RoomId    string `json:"room_id"`
	SenderId  string `json:"sender_id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
