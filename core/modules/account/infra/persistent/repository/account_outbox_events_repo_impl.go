package repos

import (
	"context"
	"go-socket/core/modules/account/domain/repos"
	"go-socket/core/modules/account/infra/persistent/models"
	"go-socket/core/shared/pkg/event"

	"gorm.io/gorm"
)

type accountOutboxEventsRepoImpl struct {
	db *gorm.DB
}

func NewAccountOutboxEventsRepoImpl(db *gorm.DB) repos.AccountOutboxEventsRepository {
	return &accountOutboxEventsRepoImpl{db: db}
}

func (a *accountOutboxEventsRepoImpl) Append(ctx context.Context, event event.Event) error {
	return a.db.Create(&models.AccountOutboxEventModel{
		EventName: event.GetName(),
		EventData: event.GetData(),
	}).Error
}
