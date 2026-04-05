package support

import (
	"context"
	"errors"

	"go-socket/core/shared/infra/xpaseto"
)

func AccountPayloadFromCtx(ctx context.Context) (*xpaseto.PasetoPayload, error) {
	payload, ok := ctx.Value("account").(*xpaseto.PasetoPayload)
	if !ok || payload == nil || payload.AccountID == "" {
		return nil, errors.New("unauthorized")
	}
	return payload, nil
}

func AccountIDFromCtx(ctx context.Context) (string, error) {
	payload, err := AccountPayloadFromCtx(ctx)
	if err != nil {
		return "", err
	}
	return payload.AccountID, nil
}
