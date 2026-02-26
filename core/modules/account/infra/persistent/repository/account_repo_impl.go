package repos

import (
	"context"

	"go-socket/core/modules/account/domain/entity"
	accountrepos "go-socket/core/modules/account/domain/repos"
	valueobject "go-socket/core/modules/account/domain/value_object"
	accountcache "go-socket/core/modules/account/infra/cache"
	"go-socket/core/modules/account/infra/persistent/models"
	sharedcache "go-socket/core/shared/infra/cache"

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

	entity, err := r.toEntity(&m)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	entity, err := r.toEntity(&m)
	if err != nil {
		return nil, err
	}
	_ = r.accountCache.SetByEmail(ctx, entity)
	return entity, nil
}

func (r *accountRepoImpl) CreateAccount(ctx context.Context, account *entity.Account) error {
	m := r.toModel(account)

	err := r.db.WithContext(ctx).
		Create(m).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *accountRepoImpl) UpdateAccount(ctx context.Context, account *entity.Account) error {
	m := r.toModel(account)

	err := r.db.WithContext(ctx).
		Save(m).Error
	if err != nil {
		return err
	}

	entity, err := r.toEntity(m)
	if err != nil {
		return err
	}
	_ = r.accountCache.Set(ctx, entity)
	_ = r.accountCache.SetByEmail(ctx, entity)

	return nil
}

func (r *accountRepoImpl) DeleteAccount(ctx context.Context, id string) error {
	if cached, ok, err := r.accountCache.Get(ctx, id); err == nil && ok {
		_ = r.accountCache.DeleteByEmail(ctx, cached.Email.Value())
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

	result := make([]*entity.Account, 0, len(accounts))

	for _, account := range accounts {
		e, err := r.toEntity(account)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}

	return result, nil
}

func (r *accountRepoImpl) toEntity(m *models.AccountModel) (*entity.Account, error) {
	email, err := valueobject.NewEmail(m.Email)
	if err != nil {
		return nil, err
	}
	password, err := valueobject.NewPassword(m.Password)
	if err != nil {
		return nil, err
	}
	return &entity.Account{
		ID:        m.ID,
		Email:     email,
		Password:  password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *accountRepoImpl) toModel(e *entity.Account) *models.AccountModel {
	return &models.AccountModel{
		ID:        e.ID,
		Email:     e.Email.Value(),
		Password:  e.Password.Value(),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
