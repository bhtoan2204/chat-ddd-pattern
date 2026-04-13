package projection

import (
	"context"
	"sort"
	"strings"
	"time"

	"go-socket/core/modules/room/application/projection"
	"go-socket/core/modules/room/domain/entity"
	roomrepos "go-socket/core/modules/room/domain/repos"
	"go-socket/core/modules/room/infra/projection/cassandra/views"
	"go-socket/core/shared/config"
	"go-socket/core/shared/pkg/stackErr"
	"go-socket/core/shared/utils"

	"github.com/gocql/gocql"
)

type queryRepoImpl struct {
	roomReadRepo       projection.RoomReadRepository
	messageReadRepo    projection.MessageReadRepository
	roomMemberReadRepo projection.RoomMemberReadRepository
}

func NewQueryRepoImpl(
	cfg config.CassandraConfig,
	session *gocql.Session,
) (projection.QueryRepos, error) {
	store, err := NewCassandraProjectionStore(cfg, session)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	var (
		roomRepo    roomrepos.RoomRepository
		messageRepo roomrepos.MessageRepository
	)

	return &queryRepoImpl{
		roomReadRepo: &roomQueryRepo{
			store:    store,
			roomRepo: roomRepo,
		},
		messageReadRepo: &messageQueryRepo{
			store:       store,
			messageRepo: messageRepo,
		},
		roomMemberReadRepo: &roomMemberQueryRepo{store: store},
	}, nil
}

func (r *queryRepoImpl) RoomReadRepository() projection.RoomReadRepository {
	return r.roomReadRepo
}

func (r *queryRepoImpl) MessageReadRepository() projection.MessageReadRepository {
	return r.messageReadRepo
}

func (r *queryRepoImpl) RoomMemberReadRepository() projection.RoomMemberReadRepository {
	return r.roomMemberReadRepo
}

type roomQueryRepo struct {
	store    *cassandraProjectionStore
	roomRepo roomrepos.RoomRepository
}

func (r *roomQueryRepo) ListRooms(ctx context.Context, options utils.QueryOptions) ([]*views.RoomView, error) {
	if r.roomRepo == nil {
		return r.store.ListRooms(ctx, options)
	}

	rooms, err := r.roomRepo.ListRooms(ctx, options)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	results := make([]*views.RoomView, 0, len(rooms))
	for _, room := range rooms {
		results = append(results, roomEntityToView(room))
	}
	return results, nil
}

func (r *roomQueryRepo) ListRoomsByAccount(ctx context.Context, accountID string, options utils.QueryOptions) ([]*views.RoomView, error) {
	return r.store.ListRoomsByAccount(ctx, accountID, options)
}

func (r *roomQueryRepo) GetRoomByID(ctx context.Context, id string) (*views.RoomView, error) {
	if r.roomRepo == nil {
		return r.store.GetRoomByID(ctx, id)
	}

	room, err := r.roomRepo.GetRoomByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return roomEntityToView(room), nil
}

type messageQueryRepo struct {
	store       *cassandraProjectionStore
	messageRepo roomrepos.MessageRepository
}

func (r *messageQueryRepo) GetMessageByID(ctx context.Context, id string) (*views.MessageView, error) {
	if r.messageRepo == nil {
		return r.store.GetMessageByID(ctx, id)
	}

	message, err := r.messageRepo.GetMessageByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return messageEntityToView(message), nil
}

func (r *messageQueryRepo) GetLastMessage(ctx context.Context, roomID string) (*views.MessageView, error) {
	return r.store.GetLastMessage(ctx, roomID)
}

func (r *messageQueryRepo) ListMessages(
	ctx context.Context,
	accountID,
	roomID string,
	options projection.MessageListOptions,
) ([]*views.MessageView, error) {
	if r.messageRepo != nil && options.BeforeAt == nil && strings.TrimSpace(options.BeforeID) != "" {
		message, err := r.messageRepo.GetMessageByID(ctx, strings.TrimSpace(options.BeforeID))
		if err != nil {
			return nil, stackErr.Error(err)
		}
		if message != nil {
			beforeAt := message.CreatedAt.UTC()
			options.BeforeAt = &beforeAt
		}
	}

	// Resolve cursor from canonical storage so Cassandra only keeps timeline-optimized tables.
	options.BeforeID = ""
	return r.store.ListMessages(ctx, accountID, roomID, options)
}

func (r *messageQueryRepo) GetMessageReceipt(ctx context.Context, messageID, accountID string) (string, *time.Time, *time.Time, error) {
	return r.store.GetMessageReceipt(ctx, messageID, accountID)
}

func (r *messageQueryRepo) CountMessageReceiptsByStatus(ctx context.Context, messageID, status string) (int64, error) {
	return r.store.CountMessageReceiptsByStatus(ctx, messageID, status)
}

func (r *messageQueryRepo) CountUnreadMessages(ctx context.Context, roomID, accountID string, lastReadAt *time.Time) (int64, error) {
	return r.store.CountUnreadMessages(ctx, roomID, accountID, lastReadAt)
}

type roomMemberQueryRepo struct {
	store *cassandraProjectionStore
}

func (r *roomMemberQueryRepo) ListRoomMembers(ctx context.Context, roomID string) ([]*views.RoomMemberView, error) {
	members, err := r.store.ListRoomMembers(ctx, roomID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return r.enrichMembers(ctx, members)
}

func (r *roomMemberQueryRepo) GetRoomMemberByAccount(ctx context.Context, roomID, accountID string) (*views.RoomMemberView, error) {
	member, err := r.store.GetRoomMemberByAccount(ctx, roomID, accountID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if member == nil {
		return nil, nil
	}

	members, err := r.enrichMembers(ctx, []*views.RoomMemberView{member})
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if len(members) == 0 {
		return nil, nil
	}
	return members[0], nil
}

func (r *roomMemberQueryRepo) SearchMentionCandidates(
	ctx context.Context,
	roomID,
	keyword,
	excludeAccountID string,
	limit int,
) ([]*views.MentionCandidateView, error) {
	members, err := r.ListRoomMembers(ctx, roomID)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	normalizedKeyword := strings.ToLower(strings.TrimSpace(keyword))
	excludeAccountID = strings.TrimSpace(excludeAccountID)

	results := make([]*views.MentionCandidateView, 0, len(members))
	for _, member := range members {
		if member == nil {
			continue
		}

		accountID := strings.TrimSpace(member.AccountID)
		if accountID == "" || accountID == excludeAccountID {
			continue
		}

		if normalizedKeyword != "" &&
			!strings.Contains(strings.ToLower(member.DisplayName), normalizedKeyword) &&
			!strings.Contains(strings.ToLower(member.Username), normalizedKeyword) &&
			!strings.Contains(strings.ToLower(accountID), normalizedKeyword) {
			continue
		}

		results = append(results, &views.MentionCandidateView{
			AccountID:       accountID,
			DisplayName:     strings.TrimSpace(member.DisplayName),
			Username:        strings.TrimSpace(member.Username),
			AvatarObjectKey: strings.TrimSpace(member.AvatarObjectKey),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		leftName := strings.ToLower(firstNonEmpty(results[i].DisplayName, results[i].AccountID))
		rightName := strings.ToLower(firstNonEmpty(results[j].DisplayName, results[j].AccountID))
		if leftName != rightName {
			return leftName < rightName
		}

		leftUsername := strings.ToLower(results[i].Username)
		rightUsername := strings.ToLower(results[j].Username)
		if leftUsername != rightUsername {
			return leftUsername < rightUsername
		}

		return results[i].AccountID < results[j].AccountID
	})

	limit = normalizeMentionLimit(limit)
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func (r *roomMemberQueryRepo) enrichMembers(ctx context.Context, members []*views.RoomMemberView) ([]*views.RoomMemberView, error) {
	return members, nil
}

func normalizeMentionLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 50 {
		return 50
	}
	return limit
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func roomEntityToView(room *entity.Room) *views.RoomView {
	if room == nil {
		return nil
	}

	return &views.RoomView{
		ID:              strings.TrimSpace(room.ID),
		Name:            strings.TrimSpace(room.Name),
		Description:     strings.TrimSpace(room.Description),
		RoomType:        strings.TrimSpace(string(room.RoomType)),
		OwnerID:         strings.TrimSpace(room.OwnerID),
		DirectKey:       stringPtr(room.DirectKey),
		PinnedMessageID: stringPtr(room.PinnedMessageID),
		CreatedAt:       room.CreatedAt.UTC(),
		UpdatedAt:       room.UpdatedAt.UTC(),
	}
}

func messageEntityToView(message *entity.MessageEntity) *views.MessageView {
	if message == nil {
		return nil
	}

	mentions := make([]views.MessageMentionView, 0, len(message.Mentions))
	for _, mention := range message.Mentions {
		mentions = append(mentions, views.MessageMentionView{
			AccountID:   strings.TrimSpace(mention.AccountID),
			DisplayName: strings.TrimSpace(mention.DisplayName),
			Username:    strings.TrimSpace(mention.Username),
		})
	}

	return &views.MessageView{
		ID:                     strings.TrimSpace(message.ID),
		RoomID:                 strings.TrimSpace(message.RoomID),
		SenderID:               strings.TrimSpace(message.SenderID),
		Message:                message.Message,
		MessageType:            strings.TrimSpace(message.MessageType),
		Mentions:               mentions,
		MentionAll:             message.MentionAll,
		ReplyToMessageID:       strings.TrimSpace(message.ReplyToMessageID),
		ForwardedFromMessageID: strings.TrimSpace(message.ForwardedFromMessageID),
		FileName:               strings.TrimSpace(message.FileName),
		FileSize:               message.FileSize,
		MimeType:               strings.TrimSpace(message.MimeType),
		ObjectKey:              strings.TrimSpace(message.ObjectKey),
		EditedAt:               cloneTime(message.EditedAt),
		DeletedForEveryoneAt:   cloneTime(message.DeletedForEveryoneAt),
		CreatedAt:              message.CreatedAt.UTC(),
	}
}

func stringPtr(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
