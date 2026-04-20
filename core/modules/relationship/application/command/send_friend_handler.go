package command

import (
	"context"
	"time"
	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/relationship/application/dto/in"
	"wechat-clone/core/modules/relationship/application/dto/out"
	"wechat-clone/core/modules/relationship/domain/aggregate"
	repos "wechat-clone/core/modules/relationship/domain/repos"
	"wechat-clone/core/modules/relationship/support"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"
	"wechat-clone/core/shared/utils"

	"github.com/google/uuid"
)

type sendFriendRequestHandler struct {
	baseRepo repos.Repos
}

func NewSendFriendRequest(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.SendFriendRequestRequest, *out.SendFriendRequestResponse] {
	return &sendFriendRequestHandler{
		baseRepo: baseRepo,
	}
}

func (u *sendFriendRequestHandler) Handle(ctx context.Context, req *in.SendFriendRequestRequest) (*out.SendFriendRequestResponse, error) {
	accountID, err := support.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	// TODO: check if the target user exists?
	// TODO: check have user send request yet?
	// TODO: check have they become friends yet?
	newFriendRequestAggregate, err := aggregate.NewFriendRequest(uuid.NewString())
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if err := u.baseRepo.WithTransaction(ctx, func(txRepos repos.Repos) error {
		if err := newFriendRequestAggregate.Create(accountID, req.TargetUserID, utils.NullableString("I wanna be friend with you <3"), time.Now()); err != nil {
			return stackErr.Error(err)
		}
		if err := txRepos.FriendRequestAggregateRepository().Save(ctx, newFriendRequestAggregate); err != nil {
			return stackErr.Error(err)
		}
		return nil
	}); err != nil {
		return nil, stackErr.Error(err)
	}
	return &out.SendFriendRequestResponse{
		RequestID:   newFriendRequestAggregate.AggregateID(),
		RequesterID: newFriendRequestAggregate.RequesterID,
		AddresseeID: newFriendRequestAggregate.AddresseeID,
		Status:      newFriendRequestAggregate.Status.String(),
	}, nil
}
