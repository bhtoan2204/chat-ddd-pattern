package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go-socket/core/infra/persistent/models"

	"github.com/redis/go-redis/v9"
)

type AccountCache struct {
	cache Cache
}

func NewAccountCache(cache Cache) *AccountCache {
	return &AccountCache{cache: cache}
}

func (a *AccountCache) Get(ctx context.Context, id string) (*models.AccountModel, bool, error) {
	if a == nil || a.cache == nil {
		return nil, false, nil
	}
	data, err := a.cache.Get(ctx, accountCacheKey(id))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, err
	}
	var m models.AccountModel
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, false, fmt.Errorf("unmarshal account cache failed: %w", err)
	}
	return &m, true, nil
}

func (a *AccountCache) Set(ctx context.Context, m *models.AccountModel) error {
	if a == nil || a.cache == nil || m == nil {
		return nil
	}
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal account cache failed: %w", err)
	}
	return a.cache.Set(ctx, accountCacheKey(m.ID), data)
}

func (a *AccountCache) Delete(ctx context.Context, id string) error {
	if a == nil || a.cache == nil {
		return nil
	}
	return a.cache.Delete(ctx, accountCacheKey(id))
}

func accountCacheKey(id string) string {
	return "account:" + id
}
