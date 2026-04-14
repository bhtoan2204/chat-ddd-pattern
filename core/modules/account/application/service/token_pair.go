package service

import (
	"context"
	"fmt"
	"time"

	"go-socket/core/modules/account/domain/entity"
	"go-socket/core/shared/infra/xpaseto"
	"go-socket/core/shared/pkg/stackErr"
)

func issueTokenPair(
	ctx context.Context,
	pasetoSvc xpaseto.PasetoService,
	account *entity.Account,
) (string, time.Time, string, time.Time, error) {
	if pasetoSvc == nil {
		return "", time.Time{}, "", time.Time{}, stackErr.Error(fmt.Errorf("paseto service is required"))
	}
	if account == nil {
		return "", time.Time{}, "", time.Time{}, stackErr.Error(fmt.Errorf("account snapshot is required"))
	}

	accessToken, accessExpiresAt, err := pasetoSvc.GenerateAccessToken(ctx, account)
	if err != nil {
		return "", time.Time{}, "", time.Time{}, stackErr.Error(fmt.Errorf("generate access token failed: %v", err))
	}

	refreshToken, refreshExpiresAt, err := pasetoSvc.GenerateRefreshToken(ctx, account)
	if err != nil {
		return "", time.Time{}, "", time.Time{}, stackErr.Error(fmt.Errorf("generate refresh token failed: %v", err))
	}

	return accessToken, accessExpiresAt, refreshToken, refreshExpiresAt, nil
}
