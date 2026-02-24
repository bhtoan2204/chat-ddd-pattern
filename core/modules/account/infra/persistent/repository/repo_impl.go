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
	return newRepoImplWithDB(appCtx, appCtx.GetDB())
}

func newRepoImplWithDB(appCtx *appCtx.AppContext, db *gorm.DB) repos.Repos {
	accountRepo := NewAccountRepoImpl(db, appCtx.GetCache())
	accountOutboxEventsRepo := NewAccountOutboxEventsRepoImpl(db)
	return &repoImpl{
		appCtx: appCtx,
		db:     db,

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

func (r *repoImpl) WithTransaction(ctx context.Context, fn func(repos.Repos) error) (err error) {
	log := logging.FromContext(ctx).Named("StartTransaction")
	tx := r.db.WithContext(ctx).Begin()
	if beginErr := tx.Error; beginErr != nil {
		log.Errorw("Failed to begin transaction", zap.Error(beginErr))
		return beginErr
	}
	tr := newRepoImplWithDB(r.appCtx, tx)

	defer func() {
		if rec := recover(); rec != nil {
			_ = tx.Rollback().Error
			log.Errorw("Panic -> rollback", zap.Any("panic", rec))
			panic(rec)
		}
		if err != nil {
			_ = tx.Rollback().Error
			log.Errorw("Transaction rollback", zap.Error(err))
			return
		}
		if commitErr := tx.Commit().Error; commitErr != nil {
			log.Errorw("Commit failed", zap.Error(commitErr))
			err = commitErr
		} else {
			log.Info("Transaction committed")
		}
	}()

	err = fn(tr)

	return err
}
