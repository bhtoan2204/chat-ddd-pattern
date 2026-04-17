package service

import (
	"context"

	"wechat-clone/core/modules/room/application/projection"
	roomsupport "wechat-clone/core/modules/room/application/support"
	apptypes "wechat-clone/core/modules/room/application/types"
	"wechat-clone/core/modules/room/infra/projection/cassandra/views"
	"wechat-clone/core/shared/pkg/stackErr"
	"wechat-clone/core/shared/utils"

	"github.com/samber/lo"
)

type RoomQueryService interface {
	GetRoom(ctx context.Context, query apptypes.GetRoomQuery) (*apptypes.RoomResult, error)
	ListRooms(ctx context.Context, query apptypes.ListRoomsQuery) (*apptypes.ListRoomsResult, error)
}

type roomQueryService struct {
	readRepos projection.QueryRepos
}

func newRoomQueryService(readRepos projection.QueryRepos) RoomQueryService {
	return &roomQueryService{readRepos: readRepos}
}

func (s *roomQueryService) GetRoom(ctx context.Context, query apptypes.GetRoomQuery) (*apptypes.RoomResult, error) {
	room, err := s.readRepos.RoomReadRepository().GetRoomByID(ctx, query.ID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return roomsupport.BuildRoomResult(room), nil
}

func (s *roomQueryService) ListRooms(ctx context.Context, query apptypes.ListRoomsQuery) (*apptypes.ListRoomsResult, error) {
	page := query.Page
	if page < 0 {
		page = 0
	}
	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}

	rooms, err := s.readRepos.RoomReadRepository().ListRooms(ctx, utils.QueryOptions{
		Offset:         &page,
		Limit:          &limit,
		OrderBy:        "updated_at",
		OrderDirection: "DESC",
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	result := &apptypes.ListRoomsResult{
		Page:  page,
		Limit: limit,
		Rooms: lo.FilterMap(rooms, func(room *views.RoomView, _ int) (apptypes.RoomResult, bool) {
			roomResult := roomsupport.BuildRoomResult(room)
			if roomResult == nil {
				return apptypes.RoomResult{}, false
			}
			return *roomResult, true
		}),
	}

	return result, nil
}
