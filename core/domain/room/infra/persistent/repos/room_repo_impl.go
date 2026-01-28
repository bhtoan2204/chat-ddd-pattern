package repository

import (
	"context"

	"go-socket/core/domain/room/entity"
	roomcache "go-socket/core/domain/room/infra/cache"
	"go-socket/core/domain/room/infra/persistent/models"
	"go-socket/core/domain/room/repos"
	sharedcache "go-socket/core/shared/infra/cache"
	"go-socket/core/shared/pkg/logging"
	"go-socket/utils"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type roomRepoImpl struct {
	db        *gorm.DB
	roomCache *roomcache.RoomCache
}

func NewRoomRepoImpl(db *gorm.DB, sharedCache sharedcache.Cache) repos.RoomRepository {
	return &roomRepoImpl{
		db:        db,
		roomCache: roomcache.NewRoomCache(sharedCache),
	}
}

func (r *roomRepoImpl) CreateRoom(ctx context.Context, room *entity.Room) error {
	m := r.toModel(room)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	_ = r.roomCache.Set(ctx, r.toEntity(m))
	return nil
}

func (r *roomRepoImpl) ListRooms(ctx context.Context, options utils.QueryOptions) ([]*entity.Room, error) {
	logger := logging.FromContext(ctx)

	var rooms []*models.RoomModel
	tx := r.db.WithContext(ctx)
	for _, c := range options.Conditions {
		tx = tx.Where(c.BuildCondition(), c.Value)
	}
	if options.Limit != nil {
		tx = tx.Limit(*options.Limit)
	}
	if options.Offset != nil {
		tx = tx.Offset(*options.Offset)
	}
	if options.OrderBy != "" && options.OrderDirection != "" {
		tx = tx.Order(options.OrderBy + " " + options.OrderDirection)
	}

	if err := tx.Find(&rooms).Error; err != nil {
		logger.Errorw("list rooms failed", zap.Error(err))
		return nil, err
	}
	return lo.Map(rooms, func(room *models.RoomModel, _ int) *entity.Room {
		return r.toEntity(room)
	}), nil
}

func (r *roomRepoImpl) GetRoomByID(ctx context.Context, id string) (*entity.Room, error) {
	if cached, ok, err := r.roomCache.Get(ctx, id); err == nil && ok {
		return cached, nil
	}
	var m models.RoomModel
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	_ = r.roomCache.Set(ctx, r.toEntity(&m))
	return r.toEntity(&m), nil
}

func (r *roomRepoImpl) UpdateRoom(ctx context.Context, room *entity.Room) error {
	m := r.toModel(room)
	if err := r.db.WithContext(ctx).Save(m).Error; err != nil {
		return err
	}
	_ = r.roomCache.Set(ctx, r.toEntity(m))
	return nil
}

func (r *roomRepoImpl) DeleteRoom(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&models.RoomModel{}, "id = ?", id).Error; err != nil {
		return err
	}
	return r.roomCache.Delete(ctx, id)
}

func (r *roomRepoImpl) toEntity(m *models.RoomModel) *entity.Room {
	return &entity.Room{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		RoomType:    m.RoomType,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *roomRepoImpl) toModel(e *entity.Room) *models.RoomModel {
	return &models.RoomModel{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		RoomType:    e.RoomType,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
