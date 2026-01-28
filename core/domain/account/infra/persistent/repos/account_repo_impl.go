package repos

import (
	"context"

	"go-socket/core/domain/account/entity"
	accountcache "go-socket/core/domain/account/infra/cache"
	"go-socket/core/domain/account/infra/persistent/models"
	accountrepos "go-socket/core/domain/account/repos"
	sharedcache "go-socket/core/shared/infra/cache"
	"go-socket/core/shared/pkg/logging"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type accountRepoImpl struct {
	db           *gorm.DB
	accountCache accountcache.AccountCache
}

func NewAccountRepoImpl(db *gorm.DB, sharedCache sharedcache.Cache) accountrepos.AccountRepository {
	return &accountRepoImpl{
		db:           db,
		accountCache: accountcache.NewAccountCache(sharedCache),
	}
}

func (r *accountRepoImpl) GetAccountByID(ctx context.Context, id string) (*entity.Account, error) {
	if cached, ok, err := r.accountCache.Get(ctx, id); err == nil && ok {
		return cached, nil
	}
	var m models.AccountModel

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&m).Error

	if err != nil {
		return nil, err
	}

	_ = r.accountCache.Set(ctx, r.toEntity(&m))

	return r.toEntity(&m), nil
}

func (r *accountRepoImpl) GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error) {
	logger := logging.FromContext(ctx)
	if cached, ok, err := r.accountCache.GetByEmail(ctx, email); err == nil && ok {
		return cached, nil
	}
	var m models.AccountModel
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&m).Error
	if err != nil {
		logger.Errorw("Failed to get account by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}
	_ = r.accountCache.SetByEmail(ctx, r.toEntity(&m))
	return r.toEntity(&m), nil
}

func (r *accountRepoImpl) CreateAccount(ctx context.Context, account *entity.Account) error {
	m := r.toModel(account)

	err := r.db.WithContext(ctx).
		Create(m).Error
	if err != nil {
		return err
	}

	_ = r.accountCache.Set(ctx, r.toEntity(m))
	_ = r.accountCache.SetByEmail(ctx, r.toEntity(m))
	return nil
}

func (r *accountRepoImpl) UpdateAccount(ctx context.Context, account *entity.Account) error {
	m := r.toModel(account)

	err := r.db.WithContext(ctx).
		Save(m).Error
	if err != nil {
		return err
	}

	_ = r.accountCache.Set(ctx, r.toEntity(m))
	_ = r.accountCache.SetByEmail(ctx, r.toEntity(m))
	return nil
}

func (r *accountRepoImpl) DeleteAccount(ctx context.Context, id string) error {
	if cached, ok, err := r.accountCache.Get(ctx, id); err == nil && ok {
		_ = r.accountCache.DeleteByEmail(ctx, cached.Email)
	}
	err := r.db.WithContext(ctx).
		Delete(&models.AccountModel{}, "id = ?", id).Error
	if err != nil {
		return err
	}

	return r.accountCache.Delete(ctx, id)
}

func (r *accountRepoImpl) ListAccountsByRoomID(ctx context.Context, roomID string) ([]*entity.Account, error) {
	var accounts []*models.AccountModel
	err := r.db.WithContext(ctx).
		Model(&models.AccountModel{}).
		Select("accounts.*").
		Joins("JOIN room_members rm ON rm.account_id = accounts.id").
		Where("rm.room_id = ?", roomID).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	return lo.Map(accounts, func(account *models.AccountModel, _ int) *entity.Account {
		return r.toEntity(account)
	}), nil
}

func (r *accountRepoImpl) toEntity(m *models.AccountModel) *entity.Account {
	return &entity.Account{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *accountRepoImpl) toModel(e *entity.Account) *models.AccountModel {
	return &models.AccountModel{
		ID:        e.ID,
		Email:     e.Email,
		Password:  e.Password,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
