package repos

import (
	"context"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/pkg/logging"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repoImpl struct {
	db     *gorm.DB
	appCtx *appCtx.AppContext

	accountRepo             repos.AccountRepository
	accountOutboxEventsRepo repos.AccountOutboxEventsRepository
}

func NewRepoImpl(appCtx *appCtx.AppContext) repos.Repos {
	accountRepo := NewAccountRepoImpl(appCtx.GetDB(), appCtx.GetCache())
	accountOutboxEventsRepo := NewAccountOutboxEventsRepoImpl(appCtx.GetDB())
	return &repoImpl{
		appCtx: appCtx,

		accountRepo:             accountRepo,
		accountOutboxEventsRepo: accountOutboxEventsRepo,
	}
}

func (r *repoImpl) AccountRepository() repos.AccountRepository {
	return r.accountRepo
}

func (r *repoImpl) AccountOutboxEventsRepository() repos.AccountOutboxEventsRepository {
	return r.accountOutboxEventsRepo
}

func (r *repoImpl) WithTransaction(ctx context.Context, fn func(repos.Repos) error) error {
	log := logging.FromContext(ctx).Named("StartTransaction")
	tx := r.db.Begin()
	if err := tx.Error; err != nil {
		log.Errorw("Failed to begin transaction", zap.Error(err))
		return err
	}
	tr := NewRepoImpl(r.appCtx)

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			log.Errorw("Failed to rollback transaction", zap.Error(err.(error)))
			panic(err)
		} else if err != nil {
			tx.Rollback()
			log.Errorw("Failed to commit transaction", zap.Error(err.(error)))
		} else {
			log.Info("Committing transaction")
			tx.Commit()
		}
	}()

	return fn(tr)
}
