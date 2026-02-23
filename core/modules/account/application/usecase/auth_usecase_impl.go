package usecase

import (
	"context"
	"errors"
	"fmt"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/domain/entity"
	accountrepos "go-socket/core/modules/account/domain/repos"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/infra/xpaseto"
	"go-socket/core/shared/pkg/hasher"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authUsecaseImpl struct {
	accountRepo accountrepos.AccountRepository
	hasher      hasher.Hasher
	paseto      xpaseto.PasetoService
}

func NewAuthUsecase(appCtx *appCtx.AppContext, repos repos.Repos) AuthUsecase {
	return &authUsecaseImpl{
		accountRepo: repos.AccountRepository(),
		hasher:      appCtx.GetHasher(),
		paseto:      appCtx.GetPaseto(),
	}
}

func (u *authUsecaseImpl) Login(ctx context.Context, in *in.LoginRequest) (*out.LoginResponse, error) {
	account, err := u.accountRepo.GetAccountByEmail(ctx, in.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccountNotFound
	}
	valid, err := u.hasher.Verify(ctx, in.Password, account.Password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, ErrInvalidCredentials
	}
	token, expiresAt, err := u.paseto.GenerateToken(ctx, account)
	if err != nil {
		return nil, err
	}
	return &out.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt.UnixMilli(),
	}, nil
}

func (u *authUsecaseImpl) Register(ctx context.Context, in *in.RegisterRequest) (*out.RegisterResponse, error) {
	_, err := u.accountRepo.GetAccountByEmail(ctx, in.Email)
	if err == nil {
		return nil, ErrAccountExists
	}
	hashedPassword, err := u.hasher.Hash(ctx, in.Password)
	if err != nil {
		return nil, err
	}
	newAccountEntity := &entity.Account{
		ID:       uuid.New().String(),
		Email:    in.Email,
		Password: hashedPassword,
	}
	if err := u.accountRepo.CreateAccount(ctx, newAccountEntity); err != nil {
		return nil, fmt.Errorf("create account failed: %w", err)
	}
	token, _, err := u.paseto.GenerateToken(ctx, newAccountEntity)
	if err != nil {
		return nil, fmt.Errorf("generate token failed: %w", err)
	}
	return &out.RegisterResponse{
		Token: token,
	}, nil
}

func (u *authUsecaseImpl) Logout(ctx context.Context, in *in.LogoutRequest) (*out.LogoutResponse, error) {
	return &out.LogoutResponse{
		Message: "Logout successful",
	}, nil
}

func (u *authUsecaseImpl) GetProfile(ctx context.Context, in *in.GetProfileRequest) (*out.GetProfileResponse, error) {
	account := ctx.Value("account")
	if account == nil {
		return nil, errors.New("account not found")
	}
	userID := account.(*xpaseto.PasetoPayload).AccountID
	accountEntity, err := u.accountRepo.GetAccountByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &out.GetProfileResponse{
		Email:     accountEntity.Email,
		CreatedAt: accountEntity.CreatedAt.Format(time.RFC3339),
		UpdatedAt: accountEntity.UpdatedAt.Format(time.RFC3339),
	}, nil
}
