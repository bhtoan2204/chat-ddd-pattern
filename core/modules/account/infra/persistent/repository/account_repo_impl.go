package repos

import (
	"context"

	"go-socket/core/modules/account/domain/entity"
	accountrepos "go-socket/core/modules/account/domain/repos"
	valueobject "go-socket/core/modules/account/domain/value_object"
	accountcache "go-socket/core/modules/account/infra/cache"
	"go-socket/core/modules/account/infra/persistent/models"
	sharedcache "go-socket/core/shared/infra/cache"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"

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
		return nil, stackErr.Error(err)
	}

	entity, err := r.toEntity(&m)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	_ = r.accountCache.Set(ctx, entity)

	return entity, nil
}

func (r *accountRepoImpl) GetAccountByEmail(ctx context.Context, email string) (*entity.Account, error) {
	if cached, ok, err := r.accountCache.GetByEmail(ctx, email); err == nil && ok {
		return cached, nil
	}
	var m models.AccountModel
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&m).Error
	if err != nil {
		return nil, stackErr.Error(err)
	}
	entity, err := r.toEntity(&m)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	_ = r.accountCache.SetByEmail(ctx, entity)
	return entity, nil
}

func (r *accountRepoImpl) IsEmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.AccountModel{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, stackErr.Error(err)
	}
	return count > 0, nil
}

func (r *accountRepoImpl) CreateAccount(ctx context.Context, account *entity.Account) error {
	m := r.toModel(account)

	if err := r.db.WithContext(ctx).
		Create(m).Error; err != nil {
		return stackErr.Error(err)
	}
	return nil
}

func (r *accountRepoImpl) UpdateAccount(ctx context.Context, account *entity.Account) error {
	m := r.toModel(account)

	if err := r.db.WithContext(ctx).
		Save(m).Error; err != nil {
		return stackErr.Error(err)
	}

	if entity, err := r.toEntity(m); err != nil {
		return stackErr.Error(err)
	} else {
		_ = r.accountCache.Set(ctx, entity)
		_ = r.accountCache.SetByEmail(ctx, entity)
	}
	return nil
}

func (r *accountRepoImpl) DeleteAccount(ctx context.Context, id string) error {
	log := logging.FromContext(ctx).Named("DeleteAccount")
	if cached, ok, err := r.accountCache.Get(ctx, id); err == nil && ok {
		_ = r.accountCache.DeleteByEmail(ctx, cached.Email.Value())
	}
	if err := r.db.WithContext(ctx).
		Delete(&models.AccountModel{}, "id = ?", id).Error; err != nil {
		return stackErr.Error(err)
	}

	if err := r.accountCache.Delete(ctx, id); err != nil {
		log.Errorw("Failed to delete account cache", zap.Error(err))
		return stackErr.Error(err)
	}
	return nil
}

func (r *accountRepoImpl) ListAccountsByRoomID(ctx context.Context, roomID string) ([]*entity.Account, error) {
	var accounts []*models.AccountModel
	if err := r.db.WithContext(ctx).
		Model(&models.AccountModel{}).
		Select("accounts.*").
		Joins("JOIN room_members rm ON rm.account_id = accounts.id").
		Where("rm.room_id = ?", roomID).
		Find(&accounts).Error; err != nil {
		return nil, stackErr.Error(err)
	}

	result := make([]*entity.Account, 0, len(accounts))

	for _, account := range accounts {
		e, err := r.toEntity(account)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		result = append(result, e)
	}

	return result, nil
}

func (r *accountRepoImpl) toEntity(m *models.AccountModel) (*entity.Account, error) {
	email, err := valueobject.NewEmail(m.Email)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	password, err := valueobject.NewPassword(m.Password)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return &entity.Account{
		ID:                m.ID,
		Email:             email,
		Password:          password,
		DisplayName:       m.DisplayName,
		Username:          m.Username,
		AvatarObjectKey:   m.AvatarObjectKey,
		Status:            m.Status,
		EmailVerifiedAt:   m.EmailVerifiedAt,
		LastLoginAt:       m.LastLoginAt,
		PasswordChangedAt: m.PasswordChangedAt,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		BannedReason:      m.BannedReason,
		BannedUntil:       m.BannedUntil,
	}, nil
}

func (r *accountRepoImpl) toModel(e *entity.Account) *models.AccountModel {
	return &models.AccountModel{
		ID:                e.ID,
		Email:             e.Email.Value(),
		Password:          e.Password.Value(),
		DisplayName:       e.DisplayName,
		Username:          e.Username,
		AvatarObjectKey:   e.AvatarObjectKey,
		Status:            e.Status,
		EmailVerifiedAt:   e.EmailVerifiedAt,
		LastLoginAt:       e.LastLoginAt,
		PasswordChangedAt: e.PasswordChangedAt,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
		BannedReason:      e.BannedReason,
		BannedUntil:       e.BannedUntil,
	}
}
