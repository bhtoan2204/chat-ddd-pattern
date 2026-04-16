package types

type MessagePayload struct {
	RoomId       string
	RecipientIds []string
	Type         string
	Payload      interface{}
}
