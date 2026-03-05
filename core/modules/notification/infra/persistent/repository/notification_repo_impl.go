package repository

import (
	"context"
	"go-socket/core/modules/notification/domain/entity"
	"go-socket/core/modules/notification/domain/repos"
	"go-socket/core/modules/notification/infra/persistent/models"

	"gorm.io/gorm"
)

type notificationRepoImpl struct {
	db *gorm.DB
}

func NewNotificationRepoImpl(db *gorm.DB) repos.NotificationRepository {
	return &notificationRepoImpl{db: db}
}

func (r *notificationRepoImpl) CreateNotification(ctx context.Context, notification *entity.NotificationEntity) error {
	m := r.toModel(notification)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *notificationRepoImpl) toModel(e *entity.NotificationEntity) *models.NotificationModel {
	return &models.NotificationModel{
		ID:        e.ID,
		AccountID: e.AccountID,
		Type:      e.Type,
		Subject:   e.Subject,
		Body:      e.Body,
		IsRead:    e.IsRead,
		ReadAt:    e.ReadAt,
		CreatedAt: e.CreatedAt,
	}
}
