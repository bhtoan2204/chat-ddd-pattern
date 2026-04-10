package repository

import (
	"context"
	"go-socket/core/modules/room/domain/entity"
	"go-socket/core/modules/room/domain/repos"
	"go-socket/core/modules/room/infra/persistent/models"
	"go-socket/core/shared/pkg/stackErr"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type roomAccountProjectionImpl struct {
	db *gorm.DB
}

func NewRoomAccountProjectionImpl(db *gorm.DB) repos.RoomAccountProjectionRepository {
	return &roomAccountProjectionImpl{db: db}
}

func (r *roomAccountProjectionImpl) ProjectAccount(ctx context.Context, account *entity.AccountEntity) error {
	model := r.toModel(account)

	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "account_id"},
			},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"display_name":      account.DisplayName,
				"username":          account.Username,
				"avatar_object_key": account.AvatarObjectKey,
				"updated_at":        account.UpdatedAt,
			}),
		}).
		Create(model).Error; err != nil {
		return stackErr.Error(err)
	}

	return nil
}

func (r *roomAccountProjectionImpl) ListByAccountIDs(ctx context.Context, accountIDs []string) ([]*entity.AccountEntity, error) {
	if len(accountIDs) == 0 {
		return []*entity.AccountEntity{}, nil
	}
	var models []models.RoomAccountProjection
	if err := r.db.WithContext(ctx).
		Where("account_id IN ?", accountIDs).
		Find(&models).Error; err != nil {
		return nil, stackErr.Error(err)
	}

	results := make([]*entity.AccountEntity, 0, len(models))
	for _, m := range models {
		model := m
		results = append(results, &entity.AccountEntity{
			AccountID:       model.AccountID,
			DisplayName:     model.DisplayName,
			Username:        model.Username,
			AvatarObjectKey: model.AvatarObjectKey,
			CreatedAt:       model.CreatedAt,
			UpdatedAt:       model.UpdatedAt,
		})
	}

	return results, nil
}

func (r *roomAccountProjectionImpl) toModel(account *entity.AccountEntity) *models.RoomAccountProjection {
	return &models.RoomAccountProjection{
		AccountID:       account.AccountID,
		DisplayName:     account.DisplayName,
		Username:        account.Username,
		AvatarObjectKey: account.AvatarObjectKey,
		CreatedAt:       account.CreatedAt,
		UpdatedAt:       account.UpdatedAt,
	}
}
