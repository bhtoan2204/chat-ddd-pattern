package idempotency

import (
	"context"
	"time"

	"go-socket/core/shared/infra/cache"
)

const keyPrefix = "idempotency:"

type RedisStore struct {
	cache cache.Cache
}

func NewRedisStore(cache cache.Cache) *RedisStore {
	return &RedisStore{cache: cache}
}

func (s *RedisStore) TryLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	if s == nil || s.cache == nil {
		return true, nil
	}
	seconds := int64(ttl.Seconds())
	if seconds <= 0 {
		seconds = 1
	}
	return s.cache.SetNX(ctx, keyPrefix+key, seconds, "locked")
}

func (s *RedisStore) MarkDone(ctx context.Context, key string, ttl time.Duration) error {
	if s == nil || s.cache == nil {
		return nil
	}
	seconds := int64(ttl.Seconds())
	if seconds <= 0 {
		seconds = 1
	}
	return s.cache.SetValWithExp(ctx, keyPrefix+key, "done", seconds)
}

func (s *RedisStore) Release(ctx context.Context, key string) error {
	if s == nil || s.cache == nil {
		return nil
	}
	return s.cache.Delete(ctx, keyPrefix+key)
}
