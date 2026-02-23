package repository

import (
	"context"
	"go-socket/core/modules/room/domain/entity"
	"go-socket/core/modules/room/domain/repos"
	"go-socket/core/modules/room/infra/persistent/models"

	"gorm.io/gorm"
)

type messageRepoImpl struct {
	db *gorm.DB
}

func NewMessageRepoImpl(db *gorm.DB) repos.MessageRepository {
	return &messageRepoImpl{db: db}
}

func (r *messageRepoImpl) CreateMessage(ctx context.Context, message *entity.MessageEntity) error {
	m := r.toModel(message)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *messageRepoImpl) toModel(e *entity.MessageEntity) *models.MessageModel {
	return &models.MessageModel{
		ID:        e.ID,
		RoomID:    e.RoomID,
		SenderID:  e.SenderID,
		Message:   e.Message,
		CreatedAt: e.CreatedAt,
	}
}

func (r *messageRepoImpl) toEntity(m *models.MessageModel) *entity.MessageEntity {
	return &entity.MessageEntity{
		ID:        m.ID,
		RoomID:    m.RoomID,
		SenderID:  m.SenderID,
		Message:   m.Message,
		CreatedAt: m.CreatedAt,
	}
}
