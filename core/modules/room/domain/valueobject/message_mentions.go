package valueobject

import (
	sharedevents "wechat-clone/core/shared/contracts/events"
)

type MessageMentions struct {
	Items      []sharedevents.RoomMessageMention
	MentionAll bool
	AccountIDs []string
}
