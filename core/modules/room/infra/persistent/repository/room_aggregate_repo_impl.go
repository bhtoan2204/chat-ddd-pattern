package repository

import (
	"context"
	"errors"
	"fmt"

	"go-socket/core/modules/room/domain/aggregate"
	"go-socket/core/modules/room/domain/repos"
	"go-socket/core/modules/room/infra/persistent/models"
	eventpkg "go-socket/core/shared/pkg/event"
	"go-socket/core/shared/pkg/stackErr"

	"gorm.io/gorm"
)

const roomOutboxAggregateType = "RoomAggregate"

type roomAggregateRepoImpl struct {
	db              *gorm.DB
	roomRepo        repos.RoomRepository
	roomMemberRepo  repos.RoomMemberRepository
	roomReadRepo    repos.RoomReadRepository
	memberReadRepo  repos.RoomMemberReadRepository
	messageRepo     repos.MessageRepository
	messageReadRepo repos.MessageReadRepository
	outboxRepo      repos.RoomOutboxEventsRepository
}

func newRoomAggregateRepoImpl(db *gorm.DB, roomRepo repos.RoomRepository, roomMemberRepo repos.RoomMemberRepository, roomReadRepo repos.RoomReadRepository, memberReadRepo repos.RoomMemberReadRepository, messageRepo repos.MessageRepository, messageReadRepo repos.MessageReadRepository, outboxRepo repos.RoomOutboxEventsRepository) repos.RoomAggregateRepository {
	return &roomAggregateRepoImpl{
		db:              db,
		roomRepo:        roomRepo,
		roomMemberRepo:  roomMemberRepo,
		roomReadRepo:    roomReadRepo,
		memberReadRepo:  memberReadRepo,
		messageRepo:     messageRepo,
		messageReadRepo: messageReadRepo,
		outboxRepo:      outboxRepo,
	}
}

func (r *roomAggregateRepoImpl) Load(ctx context.Context, roomID string) (*aggregate.RoomStateAggregate, error) {
	room, err := r.roomRepo.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	members, err := r.roomMemberRepo.ListRoomMembers(ctx, roomID)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	version, err := r.loadLatestOutboxVersion(ctx, roomID)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return aggregate.RestoreRoomStateAggregate(room, members, version)
}

func (r *roomAggregateRepoImpl) LoadByDirectKey(ctx context.Context, directKey string) (*aggregate.RoomStateAggregate, error) {
	room, err := r.roomRepo.GetRoomByDirectKey(ctx, directKey)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return r.Load(ctx, room.ID)
}

func (r *roomAggregateRepoImpl) Save(ctx context.Context, agg *aggregate.RoomStateAggregate) error {
	if agg == nil {
		return stackErr.Error(aggregate.ErrRoomAggregateNil)
	}
	if agg.IsDeleted() {
		return stackErr.Error(errors.New("deleted room aggregate must be removed via Delete"))
	}
	if !agg.HasPendingRoomWrite() {
		return nil
	}

	room := agg.Room()
	if room == nil {
		return stackErr.Error(aggregate.ErrRoomAggregateNil)
	}

	if agg.IsNew() {
		if err := r.roomRepo.CreateRoom(ctx, room); err != nil {
			return stackErr.Error(err)
		}
		if err := r.roomReadRepo.UpsertRoom(ctx, room); err != nil {
			return stackErr.Error(err)
		}
	} else {
		if err := r.roomRepo.UpdateRoom(ctx, room); err != nil {
			return stackErr.Error(err)
		}
		if err := r.roomReadRepo.UpdateRoom(ctx, room); err != nil {
			return stackErr.Error(err)
		}
	}

	for _, memberID := range agg.RemovedMemberIDs() {
		if err := r.roomMemberRepo.DeleteRoomMember(ctx, room.ID, memberID); err != nil {
			return stackErr.Error(err)
		}
		if err := r.memberReadRepo.DeleteRoomMember(ctx, room.ID, memberID); err != nil {
			return stackErr.Error(err)
		}
	}

	for _, member := range agg.PendingMemberUpserts() {
		if member == nil {
			continue
		}
		existing, err := r.roomMemberRepo.GetRoomMemberByAccount(ctx, member.RoomID, member.AccountID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return stackErr.Error(err)
			}
			if err := r.roomMemberRepo.CreateRoomMember(ctx, member); err != nil {
				return stackErr.Error(err)
			}
		} else if existing != nil {
			if err := r.roomMemberRepo.UpdateRoomMember(ctx, member); err != nil {
				return stackErr.Error(err)
			}
		}

		if err := r.memberReadRepo.UpsertRoomMember(ctx, member); err != nil {
			return stackErr.Error(err)
		}
	}

	pendingMessages := agg.PendingMessages()
	for _, message := range pendingMessages {
		if message == nil {
			continue
		}
		if err := r.messageRepo.CreateMessage(ctx, message); err != nil {
			return stackErr.Error(err)
		}
		if err := r.messageReadRepo.UpsertMessage(ctx, message); err != nil {
			return stackErr.Error(err)
		}
	}

	for _, receipt := range agg.PendingReceipts() {
		if err := r.messageReadRepo.UpsertMessageReceipt(
			ctx,
			receipt.MessageID,
			receipt.AccountID,
			receipt.Status,
			receipt.DeliveredAt,
			receipt.SeenAt,
			receipt.CreatedAt,
			receipt.UpdatedAt,
		); err != nil {
			return stackErr.Error(err)
		}
	}

	if len(pendingMessages) > 0 {
		lastMessage := pendingMessages[len(pendingMessages)-1]
		if err := r.roomReadRepo.UpdateRoomStats(ctx, room.ID, len(agg.Members()), lastMessage, room.UpdatedAt); err != nil {
			return stackErr.Error(err)
		}
	}

	nextVersion := agg.BaseVersion()
	for idx, pendingEvent := range agg.PendingOutboxEvents() {
		nextVersion++
		if err := r.outboxRepo.Append(ctx, eventpkg.Event{
			AggregateID:   room.ID,
			AggregateType: roomOutboxAggregateType,
			Version:       nextVersion,
			EventName:     pendingEvent.EventName,
			EventData:     pendingEvent.Payload,
			CreatedAt:     pendingEvent.CreatedAt.Unix(),
		}); err != nil {
			return stackErr.Error(fmt.Errorf("append room outbox event #%d failed: %v", idx, err))
		}
	}

	agg.MarkPersisted(nextVersion)
	return nil
}

func (r *roomAggregateRepoImpl) Delete(ctx context.Context, roomID string) error {
	if err := r.roomRepo.DeleteRoom(ctx, roomID); err != nil {
		return stackErr.Error(err)
	}
	return stackErr.Error(r.roomReadRepo.DeleteRoom(ctx, roomID))
}

func (r *roomAggregateRepoImpl) loadLatestOutboxVersion(ctx context.Context, roomID string) (int, error) {
	var result struct {
		Version int
	}

	err := r.db.WithContext(ctx).
		Model(&models.RoomOutboxEventModel{}).
		Select("COALESCE(MAX(version), 0) AS version").
		Where("aggregate_id = ?", roomID).
		Scan(&result).Error
	if err != nil {
		return 0, stackErr.Error(err)
	}
	return result.Version, nil
}
