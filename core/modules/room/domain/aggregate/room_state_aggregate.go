package aggregate

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go-socket/core/modules/room/domain/entity"
	roomtypes "go-socket/core/modules/room/types"
	sharedevents "go-socket/core/shared/contracts/events"
	"go-socket/core/shared/pkg/stackErr"
)

var ErrRoomAggregateNil = errors.New("room aggregate is nil")

type PendingMessageReceipt struct {
	MessageID   string
	AccountID   string
	Status      string
	DeliveredAt *time.Time
	SeenAt      *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PendingRoomOutboxEvent struct {
	EventName string
	Payload   interface{}
	CreatedAt time.Time
}

type MessageSenderIdentity struct {
	Name  string
	Email string
}

type MessageOutboxPayload struct {
	Mentions            []sharedevents.RoomMessageMention
	MentionAll          bool
	MentionedAccountIDs []string
}

type UpdateGroupDetailsParams struct {
	ActorID       string
	Name          string
	Description   string
	Now           time.Time
	SystemActorID string
}

type RoomStateAggregate struct {
	room             *entity.Room
	members          map[string]*entity.RoomMemberEntity
	memberOrder      []string
	baseVersion      int
	isNew            bool
	roomDirty        bool
	roomDeleted      bool
	pendingMessages  []*entity.MessageEntity
	pendingReceipts  []PendingMessageReceipt
	pendingOutbox    []PendingRoomOutboxEvent
	memberUpserts    map[string]*entity.RoomMemberEntity
	removedMemberIDs []string
}

func NewRoomStateAggregate(room *entity.Room, baseVersion int) (*RoomStateAggregate, error) {
	if room == nil {
		return nil, stackErr.Error(ErrRoomAggregateNil)
	}

	return &RoomStateAggregate{
		room:          room,
		members:       make(map[string]*entity.RoomMemberEntity),
		memberOrder:   []string{},
		baseVersion:   baseVersion,
		isNew:         true,
		roomDirty:     true,
		pendingOutbox: []PendingRoomOutboxEvent{},
		memberUpserts: make(map[string]*entity.RoomMemberEntity),
	}, nil
}

func RestoreRoomStateAggregate(room *entity.Room, members []*entity.RoomMemberEntity, baseVersion int) (*RoomStateAggregate, error) {
	agg, err := NewRoomStateAggregate(room, baseVersion)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	agg.isNew = false
	agg.roomDirty = false
	for _, member := range members {
		if err := agg.attachMember(member, false); err != nil {
			return nil, stackErr.Error(err)
		}
	}
	agg.memberUpserts = make(map[string]*entity.RoomMemberEntity)
	return agg, nil
}

func NewConversationRoomAggregate(
	room *entity.Room,
	members []*entity.RoomMemberEntity,
	systemActorID,
	systemMessage string,
	now time.Time,
) (*RoomStateAggregate, error) {
	agg, err := NewRoomStateAggregate(room, 0)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	agg.RecordCreated(len(members), now)

	for _, member := range members {
		if err := agg.attachMember(member, true); err != nil {
			return nil, stackErr.Error(err)
		}
	}

	if strings.TrimSpace(systemMessage) != "" {
		if _, err := agg.appendSystemMessage(systemActorID, systemMessage, now); err != nil {
			return nil, stackErr.Error(err)
		}
	}

	return agg, nil
}

func (a *RoomStateAggregate) RecordCreated(memberCount int, now time.Time) {
	a.enqueueOutbox(EventRoomCreated{}, &EventRoomCreated{
		RoomID:      a.room.ID,
		RoomType:    a.room.RoomType,
		MemberCount: memberCount,
	}, now)
}

func (a *RoomStateAggregate) Room() *entity.Room {
	return a.room
}

func (a *RoomStateAggregate) BaseVersion() int {
	return a.baseVersion
}

func (a *RoomStateAggregate) IsNew() bool {
	return a.isNew
}

func (a *RoomStateAggregate) IsDeleted() bool {
	return a.roomDeleted
}

func (a *RoomStateAggregate) Members() []*entity.RoomMemberEntity {
	results := make([]*entity.RoomMemberEntity, 0, len(a.memberOrder))
	for _, accountID := range a.memberOrder {
		member, ok := a.members[accountID]
		if !ok || member == nil {
			continue
		}
		results = append(results, member)
	}
	return results
}

func (a *RoomStateAggregate) PendingMessages() []*entity.MessageEntity {
	return append([]*entity.MessageEntity(nil), a.pendingMessages...)
}

func (a *RoomStateAggregate) PendingReceipts() []PendingMessageReceipt {
	return append([]PendingMessageReceipt(nil), a.pendingReceipts...)
}

func (a *RoomStateAggregate) PendingOutboxEvents() []PendingRoomOutboxEvent {
	return append([]PendingRoomOutboxEvent(nil), a.pendingOutbox...)
}

func (a *RoomStateAggregate) PendingMemberUpserts() []*entity.RoomMemberEntity {
	results := make([]*entity.RoomMemberEntity, 0, len(a.memberUpserts))
	for _, member := range a.memberUpserts {
		if member == nil {
			continue
		}
		results = append(results, member)
	}
	return results
}

func (a *RoomStateAggregate) RemovedMemberIDs() []string {
	return append([]string(nil), a.removedMemberIDs...)
}

func (a *RoomStateAggregate) MarkPersisted(baseVersion int) {
	a.baseVersion = baseVersion
	a.isNew = false
	a.roomDirty = false
	a.pendingMessages = nil
	a.pendingReceipts = nil
	a.pendingOutbox = nil
	a.memberUpserts = make(map[string]*entity.RoomMemberEntity)
	a.removedMemberIDs = nil
}

func (a *RoomStateAggregate) ChangeOwner(ownerID string, updatedAt time.Time) (bool, error) {
	if a == nil || a.room == nil {
		return false, stackErr.Error(ErrRoomAggregateNil)
	}
	changed, err := a.room.ChangeOwner(ownerID, updatedAt)
	if err != nil {
		return false, stackErr.Error(err)
	}
	if changed {
		a.roomDirty = true
	}
	return changed, nil
}

func (a *RoomStateAggregate) UpdateRoomDetails(name, description string, roomType roomtypes.RoomType, updatedAt time.Time) (bool, error) {
	if a == nil || a.room == nil {
		return false, stackErr.Error(ErrRoomAggregateNil)
	}

	updated, err := a.room.UpdateDetails(name, description, roomType, updatedAt)
	if err != nil {
		return false, stackErr.Error(err)
	}
	if updated {
		a.roomDirty = true
	}
	return updated, nil
}

func (a *RoomStateAggregate) UpdateGroupDetails(params UpdateGroupDetailsParams) (bool, error) {
	actor, err := a.requireMember(params.ActorID)
	if err != nil {
		return false, stackErr.Error(err)
	}
	if err := actor.CanManageGroup(a.room); err != nil {
		return false, stackErr.Error(err)
	}

	updated, err := a.room.UpdateDetails(params.Name, params.Description, "", params.Now)
	if err != nil {
		return false, stackErr.Error(err)
	}
	if !updated {
		return false, nil
	}

	a.roomDirty = true
	if _, err := a.appendSystemMessage(params.SystemActorID, fmt.Sprintf("group renamed to %s", a.room.Name), params.Now); err != nil {
		return false, stackErr.Error(err)
	}
	return true, nil
}

func (a *RoomStateAggregate) AddMember(actorID string, member *entity.RoomMemberEntity, now time.Time, systemActorID string) (bool, error) {
	actor, err := a.requireMember(actorID)
	if err != nil {
		return false, stackErr.Error(err)
	}
	if err := actor.CanManageGroup(a.room); err != nil {
		return false, stackErr.Error(err)
	}
	if member == nil {
		return false, stackErr.Error(entity.ErrRoomMemberRequired)
	}
	if _, exists := a.members[strings.TrimSpace(member.AccountID)]; exists {
		return false, nil
	}

	if err := a.attachMember(member, true); err != nil {
		return false, stackErr.Error(err)
	}
	a.room.Touch(now)
	a.roomDirty = true
	if _, err := a.appendSystemMessage(systemActorID, fmt.Sprintf("%s joined", member.AccountID), now); err != nil {
		return false, stackErr.Error(err)
	}
	return true, nil
}

func (a *RoomStateAggregate) RemoveMember(actorID, targetAccountID string, now time.Time, systemActorID string) (bool, error) {
	actor, err := a.requireMember(actorID)
	if err != nil {
		return false, stackErr.Error(err)
	}
	if err := actor.CanRemoveFrom(a.room, targetAccountID); err != nil {
		return false, stackErr.Error(err)
	}

	targetAccountID = strings.TrimSpace(targetAccountID)
	removedMember, ok := a.members[targetAccountID]
	if !ok || removedMember == nil {
		return false, stackErr.Error(entity.ErrRoomMemberRequired)
	}

	delete(a.members, targetAccountID)
	a.removedMemberIDs = append(a.removedMemberIDs, targetAccountID)
	delete(a.memberUpserts, targetAccountID)
	a.enqueueOutbox(EventRoomMemberRemoved{}, &EventRoomMemberRemoved{
		RoomID:         a.room.ID,
		MemberID:       removedMember.AccountID,
		MemberRole:     removedMember.Role,
		MemberJoinedAt: now.UTC(),
	}, now)

	a.room.Touch(now)
	a.roomDirty = true
	if _, err := a.appendSystemMessage(systemActorID, fmt.Sprintf("%s left", targetAccountID), now); err != nil {
		return false, stackErr.Error(err)
	}
	return true, nil
}

func (a *RoomStateAggregate) PinMessage(actorID, messageID string, now time.Time, systemActorID string) error {
	actor, err := a.requireMember(actorID)
	if err != nil {
		return stackErr.Error(err)
	}
	if err := actor.CanManageGroup(a.room); err != nil {
		return stackErr.Error(err)
	}
	if err := a.room.PinMessage(messageID, now); err != nil {
		return stackErr.Error(err)
	}

	a.roomDirty = true
	if _, err := a.appendSystemMessage(systemActorID, fmt.Sprintf("message %s pinned", a.room.PinnedMessageID), now); err != nil {
		return stackErr.Error(err)
	}
	return nil
}

func (a *RoomStateAggregate) SendMessage(
	messageID,
	senderID string,
	params entity.MessageParams,
	sender MessageSenderIdentity,
	outbox MessageOutboxPayload,
	now time.Time,
) (*entity.MessageEntity, error) {
	if _, err := a.requireMember(senderID); err != nil {
		return nil, stackErr.Error(err)
	}

	message, err := entity.NewMessage(messageID, a.room.ID, senderID, params, now)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	a.pendingMessages = append(a.pendingMessages, message)
	for _, member := range a.Members() {
		if member == nil || strings.TrimSpace(member.AccountID) == strings.TrimSpace(senderID) {
			continue
		}
		a.pendingReceipts = append(a.pendingReceipts, PendingMessageReceipt{
			MessageID: message.ID,
			AccountID: member.AccountID,
			Status:    "sent",
			CreatedAt: now.UTC(),
			UpdatedAt: now.UTC(),
		})
	}

	a.room.Touch(now)
	a.roomDirty = true
	a.enqueueOutbox(sharedevents.EventRoomMessageCreated, &sharedevents.RoomMessageCreatedEvent{
		RoomID:                 a.room.ID,
		RoomName:               a.room.Name,
		RoomType:               string(a.room.RoomType),
		MessageID:              message.ID,
		MessageContent:         message.Message,
		MessageType:            message.MessageType,
		ReplyToMessageID:       message.ReplyToMessageID,
		ForwardedFromMessageID: message.ForwardedFromMessageID,
		FileName:               message.FileName,
		FileSize:               message.FileSize,
		MimeType:               message.MimeType,
		ObjectKey:              message.ObjectKey,
		MessageSenderID:        senderID,
		MessageSenderName:      strings.TrimSpace(sender.Name),
		MessageSenderEmail:     strings.TrimSpace(sender.Email),
		MessageSentAt:          message.CreatedAt,
		Mentions:               outbox.Mentions,
		MentionAll:             outbox.MentionAll,
		MentionedAccountIDs:    outbox.MentionedAccountIDs,
	}, now)

	return message, nil
}

func (a *RoomStateAggregate) appendSystemMessage(actorID, body string, now time.Time) (*entity.MessageEntity, error) {
	message, err := entity.NewSystemMessage(newUUID(), a.room.ID, actorID, body, now)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	a.pendingMessages = append(a.pendingMessages, message)
	a.room.Touch(now)
	a.roomDirty = true
	return message, nil
}

func (a *RoomStateAggregate) HasPendingRoomWrite() bool {
	return a.roomDeleted || a.roomDirty || len(a.memberUpserts) > 0 || len(a.removedMemberIDs) > 0 || len(a.pendingMessages) > 0 || len(a.pendingReceipts) > 0 || len(a.pendingOutbox) > 0
}

func (a *RoomStateAggregate) MarkDeleted() error {
	if a == nil || a.room == nil {
		return stackErr.Error(ErrRoomAggregateNil)
	}
	a.roomDeleted = true
	return nil
}

func (a *RoomStateAggregate) attachMember(member *entity.RoomMemberEntity, emitEvent bool) error {
	if member == nil {
		return stackErr.Error(entity.ErrRoomMemberRequired)
	}

	accountID := strings.TrimSpace(member.AccountID)
	if accountID == "" {
		return stackErr.Error(entity.ErrRoomMemberAccountRequired)
	}
	if _, exists := a.members[accountID]; !exists {
		a.memberOrder = append(a.memberOrder, accountID)
	}

	a.members[accountID] = member
	a.memberUpserts[accountID] = member
	if emitEvent {
		a.enqueueOutbox(EventRoomMemberAdded{}, &EventRoomMemberAdded{
			RoomID:         a.room.ID,
			MemberID:       member.AccountID,
			MemberRole:     member.Role,
			MemberJoinedAt: member.CreatedAt,
		}, member.CreatedAt)
	}
	return nil
}

func (a *RoomStateAggregate) requireMember(accountID string) (*entity.RoomMemberEntity, error) {
	if a == nil || a.room == nil {
		return nil, stackErr.Error(ErrRoomAggregateNil)
	}

	member, ok := a.members[strings.TrimSpace(accountID)]
	if !ok || member == nil {
		return nil, stackErr.Error(entity.ErrRoomMemberRequired)
	}
	return member, nil
}

func (a *RoomStateAggregate) enqueueOutbox(eventNameOrType interface{}, payload interface{}, createdAt time.Time) {
	eventName := ""
	switch value := eventNameOrType.(type) {
	case string:
		eventName = strings.TrimSpace(value)
	default:
		eventName = strings.TrimPrefix(fmt.Sprintf("%T", value), "*aggregate.")
		if idx := strings.LastIndex(eventName, "."); idx >= 0 {
			eventName = eventName[idx+1:]
		}
	}

	a.pendingOutbox = append(a.pendingOutbox, PendingRoomOutboxEvent{
		EventName: eventName,
		Payload:   payload,
		CreatedAt: createdAt.UTC(),
	})
}
