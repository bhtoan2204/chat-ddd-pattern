package command

import (
	"context"
	"errors"
	"fmt"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/infra/xpaseto"
	"go-socket/core/shared/pkg/hasher"
	"go-socket/core/shared/pkg/logging"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type loginHandler struct {
	accountRepo repos.AccountRepository
	hasher      hasher.Hasher
	paseto      xpaseto.PasetoService
}

func NewLoginHandler(appCtx *appCtx.AppContext, baseRepo repos.Repos) LoginHandler {
	return &loginHandler{
		accountRepo: baseRepo.AccountRepository(),
		hasher:      appCtx.GetHasher(),
		paseto:      appCtx.GetPaseto(),
	}
}

func (u *loginHandler) Handle(ctx context.Context, req *in.LoginRequest) (*out.LoginResponse, error) {
	log := logging.FromContext(ctx).Named("Login")
	account, err := u.accountRepo.GetAccountByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorw("Account not found", zap.String("email", req.Email))
			return nil, ErrAccountNotFound
		}
		log.Errorw("Failed to get account", zap.Error(err))
		return nil, fmt.Errorf("get account failed: %w", err)
	}

	valid, err := u.hasher.Verify(ctx, req.Password, account.Password)
	if err != nil {
		log.Errorw("Failed to verify password", zap.Error(err))
		return nil, err
	}
	if !valid {
		log.Errorw("Invalid credentials", zap.String("email", req.Email))
		return nil, ErrInvalidCredentials
	}

	token, expiresAt, err := u.paseto.GenerateToken(ctx, account)
	if err != nil {
		log.Errorw("Failed to generate token", zap.Error(err))
		return nil, err
	}

	return &out.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt.UnixMilli(),
	}, nil
}
