package appCtx

import (
	"context"
	"go-socket/core/infra/cache"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Option func(*AppContext)

type AppContext struct {
	redisClient *redis.Client
	db          *gorm.DB
	cache       cache.Cache
}

func NewAppContext(ctx context.Context, opts ...Option) (*AppContext, error) {
	appCtx := &AppContext{}
	for _, opt := range opts {
		opt(appCtx)
	}
	return appCtx, nil
}

func WithRedisClient(redisClient *redis.Client) Option {
	return func(appCtx *AppContext) {
		appCtx.redisClient = redisClient
	}
}

func WithCache(cache cache.Cache) Option {
	return func(appCtx *AppContext) {
		appCtx.cache = cache
	}
}

func WithDB(db *gorm.DB) Option {
	return func(appCtx *AppContext) {
		appCtx.db = db
	}
}

func (appCtx *AppContext) GetRedisClient() *redis.Client {
	return appCtx.redisClient
}

func (appCtx *AppContext) GetDB() *gorm.DB {
	return appCtx.db
}

func (appCtx *AppContext) GetCache() cache.Cache {
	return appCtx.cache
}

func (appCtx *AppContext) Close() {
	appCtx.redisClient.Close()
	if appCtx.db != nil {
		ins, _ := appCtx.db.DB()
		ins.Close()
	}
}
