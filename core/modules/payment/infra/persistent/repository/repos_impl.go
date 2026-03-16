package repository

import (
	"context"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/payment/domain/repos"
	"go-socket/core/shared/pkg/logging"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repoImpl struct {
	appCtx *appCtx.AppContext
	db     *gorm.DB

	paymentBalanceAggregateRepo  repos.PaymentBalanceAggregateRepository
	paymentProjectionRepo        repos.PaymentProjectionRepository
	paymentOutboxEventsRepo      repos.PaymentOutboxEventsRepository
	paymentAccountProjectionRepo repos.PaymentAccountProjectionRepository
	paymentHistoryRepo           repos.PaymentHistoryRepository
}

func NewRepoImpl(appCtx *appCtx.AppContext) repos.Repos {
	return newRepoImplWithDB(appCtx, appCtx.GetDB())
}

func newRepoImplWithDB(appCtx *appCtx.AppContext, db *gorm.DB) repos.Repos {
	paymentBalanceAggregateRepo := NewPaymentBalanceAggregateRepoImpl(db)
	paymentProjectionRepo := NewPaymentProjectionRepoImpl(db)
	paymentOutboxEventsRepo := NewPaymentOutboxEventsRepoImpl(db)
	paymentAccountProjectionRepo := NewPaymentAccountProjectionRepoImpl(db)
	paymentHistoryRepo := NewPaymentHistoryRepoImpl(db)
	return &repoImpl{
		appCtx: appCtx,
		db:     db,

		paymentBalanceAggregateRepo:  paymentBalanceAggregateRepo,
		paymentProjectionRepo:        paymentProjectionRepo,
		paymentOutboxEventsRepo:      paymentOutboxEventsRepo,
		paymentAccountProjectionRepo: paymentAccountProjectionRepo,
		paymentHistoryRepo:           paymentHistoryRepo,
	}
}

func (r *repoImpl) PaymentBalanceAggregateRepository() repos.PaymentBalanceAggregateRepository {
	return r.paymentBalanceAggregateRepo
}

func (r *repoImpl) PaymentProjectionRepository() repos.PaymentProjectionRepository {
	return r.paymentProjectionRepo
}

func (r *repoImpl) PaymentOutboxEventsRepository() repos.PaymentOutboxEventsRepository {
	return r.paymentOutboxEventsRepo
}

func (r *repoImpl) PaymentAccountProjectionRepository() repos.PaymentAccountProjectionRepository {
	return r.paymentAccountProjectionRepo
}

func (r *repoImpl) PaymentHistoryRepository() repos.PaymentHistoryRepository {
	return r.paymentHistoryRepo
}

func (r *repoImpl) WithTransaction(ctx context.Context, fn func(repos.Repos) error) (err error) {
	log := logging.FromContext(ctx).Named("StartPaymentTransaction")
	tx := r.db.WithContext(ctx).Begin()
	if beginErr := tx.Error; beginErr != nil {
		log.Errorw("failed to begin transaction", zap.Error(beginErr))
		return beginErr
	}

	tr := newRepoImplWithDB(r.appCtx, tx)

	defer func() {
		if rec := recover(); rec != nil {
			_ = tx.Rollback().Error
			log.Errorw("panic -> rollback", zap.Any("panic", rec))
			panic(rec)
		}

		if err != nil {
			_ = tx.Rollback().Error
			log.Errorw("transaction rollback", zap.Error(err))
			return
		}

		if commitErr := tx.Commit().Error; commitErr != nil {
			log.Errorw("commit failed", zap.Error(commitErr))
			err = commitErr
		} else {
			log.Info("transaction committed")
		}
	}()

	err = fn(tr)
	return err
}
