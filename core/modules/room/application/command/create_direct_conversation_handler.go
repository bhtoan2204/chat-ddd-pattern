package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	roomsupport "go-socket/core/modules/room/application/support"
	"go-socket/core/modules/room/domain/aggregate"
	"go-socket/core/modules/room/domain/entity"
	roomrepos "go-socket/core/modules/room/domain/repos"
	roomtypes "go-socket/core/modules/room/types"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/stackErr"

	"gorm.io/gorm"
)

type createDirectConversationHandler struct {
	baseRepo roomrepos.Repos
}

func NewCreateDirectConversationHandler(baseRepo roomrepos.Repos) cqrs.Handler[*in.CreateDirectConversationRequest, *out.ChatConversationResponse] {
	return &createDirectConversationHandler{baseRepo: baseRepo}
}

func (h *createDirectConversationHandler) Handle(ctx context.Context, req *in.CreateDirectConversationRequest) (*out.ChatConversationResponse, error) {
	accountID, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if err := ensureProjectedAccountsExist(ctx, h.baseRepo, req.PeerAccountID); err != nil {
		return nil, stackErr.Error(err)
	}

	now := time.Now().UTC()
	room, err := entity.NewDirectConversationRoom(newID(), accountID, req.PeerAccountID, now)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	existing, err := h.baseRepo.RoomAggregateRepository().LoadByDirectKey(ctx, room.DirectKey)
	if err == nil && existing != nil {
		res, buildErr := roomsupport.BuildConversationResult(ctx, h.baseRepo, accountID, existing.Room(), true)
		if buildErr != nil {
			return nil, stackErr.Error(buildErr)
		}
		return roomsupport.ToConversationResponse(res), nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, stackErr.Error(err)
	}

	ownerMember, err := entity.NewRoomMember(newID(), room.ID, accountID, roomtypes.RoomRoleOwner, now)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	peerMember, err := entity.NewRoomMember(newID(), room.ID, req.PeerAccountID, roomtypes.RoomRoleMember, now)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	agg, err := aggregate.NewConversationRoomAggregate(
		room,
		[]*entity.RoomMemberEntity{ownerMember, peerMember},
		accountID,
		fmt.Sprintf("%s started a direct conversation", accountID),
		now,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	if err := h.baseRepo.WithTransaction(ctx, func(txRepos roomrepos.Repos) error {
		return stackErr.Error(txRepos.RoomAggregateRepository().Save(ctx, agg))
	}); err != nil {
		return nil, stackErr.Error(err)
	}

	res, err := roomsupport.BuildConversationResult(ctx, h.baseRepo, accountID, room, true)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return roomsupport.ToConversationResponse(res), nil
}
