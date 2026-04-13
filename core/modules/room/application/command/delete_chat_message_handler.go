package command

import (
	"context"
	"time"

	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	roomsupport "go-socket/core/modules/room/application/support"
	roomrepos "go-socket/core/modules/room/domain/repos"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/stackErr"
)

type deleteChatMessageHandler struct {
	baseRepo roomrepos.Repos
}

func NewDeleteChatMessageHandler(baseRepo roomrepos.Repos) cqrs.Handler[*in.DeleteChatMessageRequest, *out.DeleteChatMessageResponse] {
	return &deleteChatMessageHandler{baseRepo: baseRepo}
}

func (h *deleteChatMessageHandler) Handle(ctx context.Context, req *in.DeleteChatMessageRequest) (*out.DeleteChatMessageResponse, error) {
	accountID, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	agg, err := h.baseRepo.MessageAggregateRepository().Load(ctx, req.MessageID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if err := agg.Delete(accountID, accountID, req.Scope, time.Now().UTC()); err != nil {
		return nil, stackErr.Error(err)
	}

	if err := h.baseRepo.WithTransaction(ctx, func(txRepos roomrepos.Repos) error {
		return stackErr.Error(txRepos.MessageAggregateRepository().Save(ctx, agg))
	}); err != nil {
		return nil, stackErr.Error(err)
	}

	return &out.DeleteChatMessageResponse{Ok: true}, nil
}
