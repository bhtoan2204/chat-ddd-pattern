package support

import (
	"context"

	"wechat-clone/core/shared/pkg/actorctx"
)

func AccountIDFromCtx(ctx context.Context) (string, error) {
	return actorctx.AccountIDFromContext(ctx)
}
