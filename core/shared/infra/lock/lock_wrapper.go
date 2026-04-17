package lock

import (
	"context"
	"fmt"
	"sort"
	"time"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MultiLockOptions struct {
	Expiration time.Duration
	RetryDelay time.Duration
	Timeout    time.Duration
	KeyPrefix  string
}

func DefaultMultiLockOptions() MultiLockOptions {
	return MultiLockOptions{
		Expiration: 30 * time.Second,
		RetryDelay: 100 * time.Millisecond,
		Timeout:    3 * time.Second,
	}
}

func normalizeKeys(keys []string, prefix string) []string {
	if len(keys) == 0 {
		return nil
	}

	result := make([]string, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))

	for _, k := range keys {
		if k == "" {
			continue
		}

		if prefix != "" {
			k = fmt.Sprintf("%s:%s", prefix, k)
		}

		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		result = append(result, k)
	}

	sort.Strings(result)
	return result
}

func WithLocks[T any](
	ctx context.Context,
	locker Lock,
	keys []string,
	opts MultiLockOptions,
	fn func() (T, error),
) (T, error) {
	var zero T

	if locker == nil || len(keys) == 0 {
		return fn()
	}

	lockKeys := normalizeKeys(keys, opts.KeyPrefix)
	if len(lockKeys) == 0 {
		return fn()
	}

	lockValue := uuid.NewString()
	releaseKeys := make([]string, 0, len(lockKeys))

	release := func() {
		for i := len(releaseKeys) - 1; i >= 0; i-- {
			if _, err := locker.ReleaseLock(ctx, releaseKeys[i], lockValue); err != nil {
				logging.FromContext(ctx).Warnw(
					"release lock failed",
					zap.String("lock_key", releaseKeys[i]),
					zap.Error(err),
				)
			}
		}
	}

	for _, key := range lockKeys {
		locked, err := locker.AcquireLock(
			ctx,
			key,
			lockValue,
			opts.Expiration,
			opts.RetryDelay,
			opts.Timeout,
		)
		if err != nil {
			release()
			return zero, stackErr.Error(err)
		}
		if !locked {
			release()
			return zero, stackErr.Error(fmt.Errorf("acquire lock failed: %s", key))
		}

		releaseKeys = append(releaseKeys, key)
	}

	defer release()

	return fn()
}
