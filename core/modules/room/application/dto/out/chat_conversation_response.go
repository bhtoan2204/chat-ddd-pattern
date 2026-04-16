// CODE_GENERATOR - do not edit: response
package out

type ChatConversationResponse struct {
	RoomID          string                   `json:"room_id,omitempty"`
	Name            string                   `json:"name,omitempty"`
	Description     string                   `json:"description,omitempty"`
	RoomType        string                   `json:"room_type,omitempty"`
	OwnerID         string                   `json:"owner_id,omitempty"`
	PinnedMessageID string                   `json:"pinned_message_id,omitempty"`
	MemberCount     int                      `json:"member_count,omitempty"`
	UnreadCount     int64                    `json:"unread_count,omitempty"`
	LastMessage     *ChatMessageResponse     `json:"last_message,omitempty"`
	Members         []ChatRoomMemberResponse `json:"members,omitempty"`
	CreatedAt       string                   `json:"created_at,omitempty"`
	UpdatedAt       string                   `json:"updated_at,omitempty"`
}

type ChatRoomMemberResponse struct {
	AccountID       string `json:"account_id,omitempty"`
	Role            string `json:"role,omitempty"`
	DisplayName     string `json:"display_name,omitempty"`
	AvatarObjectKey string `json:"avatar_object_key,omitempty"`
}
