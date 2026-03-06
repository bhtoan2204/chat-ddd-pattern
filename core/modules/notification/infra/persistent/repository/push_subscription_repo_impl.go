package repository

import (
	"context"
	"errors"

	"go-socket/core/modules/notification/domain/entity"
	notificationrepos "go-socket/core/modules/notification/domain/repos"
	"go-socket/core/modules/notification/infra/persistent/models"
	stackerr "go-socket/core/shared/pkg/stackErr"

	"gorm.io/gorm"
)

type pushSubscriptionRepoImpl struct {
	db *gorm.DB
}

func NewPushSubscriptionRepoImpl(db *gorm.DB) notificationrepos.PushSubscriptionRepository {
	return &pushSubscriptionRepoImpl{db: db}
}

func (r *pushSubscriptionRepoImpl) UpsertPushSubscription(ctx context.Context, subscription *entity.PushSubscription) error {
	m := r.toPushSubscriptionModel(subscription)

	var existing models.PushSubscriptionModel
	err := r.db.WithContext(ctx).
		Where("account_id = ? AND endpoint = ?", subscription.AccountID, subscription.Endpoint).
		First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if createErr := r.db.WithContext(ctx).Create(m).Error; createErr != nil {
				return stackerr.Error(createErr)
			}
			return nil
		}
		return stackerr.Error(err)
	}

	if updateErr := r.db.WithContext(ctx).
		Model(&existing).
		Updates(map[string]interface{}{"keys": subscription.Keys}).Error; updateErr != nil {
		return stackerr.Error(updateErr)
	}

	return nil
}

func (r *pushSubscriptionRepoImpl) ListPushSubscriptionsByAccountID(ctx context.Context, accountID string) ([]*entity.PushSubscription, error) {
	var subscriptions []*models.PushSubscriptionModel
	if err := r.db.WithContext(ctx).
		Where("account_id = ?", accountID).
		Order("created_at DESC").
		Find(&subscriptions).Error; err != nil {
		return nil, stackerr.Error(err)
	}

	result := make([]*entity.PushSubscription, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		result = append(result, r.toPushSubscriptionEntity(subscription))
	}

	return result, nil
}

func (r *pushSubscriptionRepoImpl) toPushSubscriptionEntity(m *models.PushSubscriptionModel) *entity.PushSubscription {
	return &entity.PushSubscription{
		ID:        m.ID,
		AccountID: m.AccountID,
		Endpoint:  m.Endpoint,
		Keys:      m.Keys,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *pushSubscriptionRepoImpl) toPushSubscriptionModel(e *entity.PushSubscription) *models.PushSubscriptionModel {
	return &models.PushSubscriptionModel{
		ID:        e.ID,
		AccountID: e.AccountID,
		Endpoint:  e.Endpoint,
		Keys:      e.Keys,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
