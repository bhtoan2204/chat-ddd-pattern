package repository

import (
	"context"
	"go-socket/core/modules/room/domain/entity"
	"go-socket/core/modules/room/domain/repos"
	"go-socket/core/modules/room/infra/persistent/models"

	"gorm.io/gorm"
)

type roomMemberImpl struct {
	db *gorm.DB
}

func NewRoomMemberImpl(db *gorm.DB) repos.RoomMemberRepository {
	return &roomMemberImpl{db: db}
}

func (r *roomMemberImpl) CreateRoomMember(ctx context.Context, roomMember *entity.RoomMemberEntity) error {
	m := r.toModel(roomMember)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *roomMemberImpl) toModel(e *entity.RoomMemberEntity) *models.RoomMemberModel {
	return &models.RoomMemberModel{
		ID:        e.ID,
		RoomID:    e.RoomID,
		AccountID: e.AccountID,
		Role:      e.Role,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
