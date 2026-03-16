package command

import (
	"context"

	"go-socket/core/shared/infra/xpaseto"
)

func accountIDFromContext(ctx context.Context) (string, error) {
	account, ok := ctx.Value("account").(*xpaseto.PasetoPayload)
	if !ok || account == nil || account.AccountID == "" {
		return "", ErrPaymentAccountNotFound
	}

	return account.AccountID, nil
}
