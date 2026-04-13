package command

import (
	"context"

	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	roomrepos "go-socket/core/modules/room/domain/repos"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/stackErr"
)

type deleteRoomHandler struct {
	baseRepo roomrepos.Repos
}

func NewDeleteRoomHandler(baseRepo roomrepos.Repos) cqrs.Handler[*in.DeleteRoomRequest, *out.DeleteRoomResponse] {
	return &deleteRoomHandler{
		baseRepo: baseRepo,
	}
}

func (h *deleteRoomHandler) Handle(ctx context.Context, req *in.DeleteRoomRequest) (*out.DeleteRoomResponse, error) {
	if err := h.baseRepo.WithTransaction(ctx, func(txRepos roomrepos.Repos) error {
		return stackErr.Error(txRepos.RoomAggregateRepository().Delete(ctx, req.ID))
	}); err != nil {
		return nil, stackErr.Error(err)
	}
	return &out.DeleteRoomResponse{
		Message: "Room deleted successfully",
	}, nil
}
