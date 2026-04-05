package service

import (
	"context"
	"errors"
	"strings"
	"time"

	apptypes "go-socket/core/modules/room/application/types"
	"go-socket/core/modules/room/domain/entity"
	"go-socket/core/modules/room/domain/repos"
)

func (s *RoomCommandService) CreateRoom(ctx context.Context, accountID string, command apptypes.CreateRoomCommand) (*apptypes.RoomResult, error) {
	now := time.Now().UTC()
	room, err := entity.NewRoom(newUUID(), command.Name, command.Description, accountID, command.RoomType, "", now)
	if err != nil {
		return nil, err
	}

	if err := s.repos.WithTransaction(ctx, func(txRepos repos.Repos) error {
		if err := txRepos.RoomRepository().CreateRoom(ctx, room); err != nil {
			return err
		}
		if err := txRepos.RoomReadRepository().UpsertRoom(ctx, room); err != nil {
			return err
		}
		return s.aggregateService.PublishRoomCreated(ctx, txRepos.RoomOutboxEventsRepository(), room.ID, room.RoomType, 1)
	}); err != nil {
		return nil, err
	}

	return buildRoomResult(room), nil
}

func (s *RoomCommandService) UpdateRoom(ctx context.Context, accountID, roomID string, command apptypes.UpdateRoomCommand) (*apptypes.RoomResult, error) {
	room, err := s.repos.RoomRepository().GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	if accountID = strings.TrimSpace(accountID); accountID != "" {
		room.OwnerID = accountID
	}

	if _, err := room.UpdateDetails(command.Name, command.Description, command.RoomType, time.Now().UTC()); err != nil {
		return nil, err
	}
	if err := s.repos.WithTransaction(ctx, func(txRepos repos.Repos) error {
		if err := txRepos.RoomRepository().UpdateRoom(ctx, room); err != nil {
			return err
		}
		return txRepos.RoomReadRepository().UpdateRoom(ctx, room)
	}); err != nil {
		return nil, err
	}

	return buildRoomResult(room), nil
}

func (s *RoomCommandService) DeleteRoom(ctx context.Context, roomID string) error {
	return s.repos.WithTransaction(ctx, func(txRepos repos.Repos) error {
		if err := txRepos.RoomRepository().DeleteRoom(ctx, roomID); err != nil {
			return err
		}
		return txRepos.RoomReadRepository().DeleteRoom(ctx, roomID)
	})
}

func (s *RoomCommandService) JoinRoom(ctx context.Context, accountID string, command apptypes.JoinRoomCommand) error {
	return errors.New("not implemented")
}
