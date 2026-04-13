package repository

import (
	"context"

	"go-socket/core/modules/room/domain/aggregate"
	"go-socket/core/modules/room/domain/repos"
	"go-socket/core/shared/pkg/stackErr"

	"gorm.io/gorm"
)

type messageAggregateRepoImpl struct {
	db              *gorm.DB
	messageRepo     repos.MessageRepository
	messageReadRepo repos.MessageReadRepository
	roomMemberRepo  repos.RoomMemberRepository
	memberReadRepo  repos.RoomMemberReadRepository
}

func newMessageAggregateRepoImpl(db *gorm.DB, messageRepo repos.MessageRepository, messageReadRepo repos.MessageReadRepository, roomMemberRepo repos.RoomMemberRepository, memberReadRepo repos.RoomMemberReadRepository) repos.MessageAggregateRepository {
	return &messageAggregateRepoImpl{
		db:              db,
		messageRepo:     messageRepo,
		messageReadRepo: messageReadRepo,
		roomMemberRepo:  roomMemberRepo,
		memberReadRepo:  memberReadRepo,
	}
}

func (r *messageAggregateRepoImpl) Load(ctx context.Context, messageID string) (*aggregate.MessageStateAggregate, error) {
	message, err := r.messageRepo.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return aggregate.NewMessageStateAggregate(message)
}

func (r *messageAggregateRepoImpl) Save(ctx context.Context, agg *aggregate.MessageStateAggregate) error {
	if agg == nil {
		return stackErr.Error(aggregate.ErrMessageAggregateNil)
	}

	if agg.MessageDirty() {
		if err := r.messageRepo.UpdateMessage(ctx, agg.Message()); err != nil {
			return stackErr.Error(err)
		}
		if err := r.messageReadRepo.UpsertMessage(ctx, agg.Message()); err != nil {
			return stackErr.Error(err)
		}
	}

	if deletion := agg.PendingDeletion(); deletion != nil {
		if err := r.messageReadRepo.UpsertMessageDeletion(ctx, deletion.MessageID, deletion.AccountID, deletion.CreatedAt); err != nil {
			return stackErr.Error(err)
		}
	}

	if receipt := agg.PendingReceipt(); receipt != nil {
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

	if agg.MemberDirty() && agg.RecipientMember() != nil {
		if err := r.roomMemberRepo.UpdateRoomMember(ctx, agg.RecipientMember()); err != nil {
			return stackErr.Error(err)
		}
		if err := r.memberReadRepo.UpsertRoomMember(ctx, agg.RecipientMember()); err != nil {
			return stackErr.Error(err)
		}
	}

	agg.MarkPersisted()
	return nil
}
