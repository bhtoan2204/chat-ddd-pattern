package usecase

import (
	"context"
	"errors"
	"fmt"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/domain/entity"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/contracts/events"
	"go-socket/core/shared/infra/xpaseto"
	"go-socket/core/shared/pkg/hasher"
	"go-socket/core/shared/pkg/logging"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type authUsecaseImpl struct {
	baseRepo repos.Repos
	hasher   hasher.Hasher
	paseto   xpaseto.PasetoService
}

func NewAuthUsecase(appCtx *appCtx.AppContext, repos repos.Repos) AuthUsecase {
	return &authUsecaseImpl{
		baseRepo: repos,
		hasher:   appCtx.GetHasher(),
		paseto:   appCtx.GetPaseto(),
	}
}

func (u *authUsecaseImpl) Login(ctx context.Context, in *in.LoginRequest) (*out.LoginResponse, error) {
	log := logging.FromContext(ctx).Named("Login")
	accountRepo := u.baseRepo.AccountRepository()
	account, err := accountRepo.GetAccountByEmail(ctx, in.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorw("Account not found", zap.String("email", in.Email))
		return nil, ErrAccountNotFound
	}
	valid, err := u.hasher.Verify(ctx, in.Password, account.Password)
	if err != nil {
		log.Errorw("Failed to verify password", zap.Error(err))
		return nil, err
	}
	if !valid {
		log.Errorw("Invalid credentials", zap.String("email", in.Email))
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

func (u *authUsecaseImpl) Register(ctx context.Context, in *in.RegisterRequest) (*out.RegisterResponse, error) {
	log := logging.FromContext(ctx).Named("Register")
	_, err := u.baseRepo.AccountRepository().GetAccountByEmail(ctx, in.Email)
	if err == nil {
		log.Errorw("Account already exists", zap.String("email", in.Email))
		return nil, ErrAccountExists
	}
	hashedPassword, err := u.hasher.Hash(ctx, in.Password)
	if err != nil {
		log.Errorw("Failed to hash password", zap.Error(err))
		return nil, err
	}
	newAccountEntity := &entity.Account{
		ID:       uuid.New().String(),
		Email:    in.Email,
		Password: hashedPassword,
	}
	if txErr := u.baseRepo.WithTransaction(ctx, func(repos repos.Repos) error {
		if err := repos.AccountRepository().CreateAccount(ctx, newAccountEntity); err != nil {
			log.Errorw("Failed to create account", zap.Error(err))
			return fmt.Errorf("create account failed: %w", err)
		}
		if err := repos.AccountOutboxEventsRepository().Append(ctx, &events.AccountCreatedEvent{
			AccountID: newAccountEntity.ID,
			Email:     newAccountEntity.Email,
			CreatedAt: time.Now(),
		}); err != nil {
			log.Errorw("Failed to append account created event", zap.Error(err))
			return fmt.Errorf("append account created event failed: %w", err)
		}
		return nil
	}); txErr != nil {
		log.Errorw("Failed to register account", zap.Error(txErr))
		return nil, txErr
	}

	token, _, err := u.paseto.GenerateToken(ctx, newAccountEntity)
	if err != nil {
		log.Errorw("Failed to generate token", zap.Error(err))
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
	log := logging.FromContext(ctx).Named("GetProfile")
	account := ctx.Value("account")
	if account == nil {
		log.Errorw("Account not found", zap.Error(errors.New("account not found")))
		return nil, errors.New("account not found")
	}
	accountRepo := u.baseRepo.AccountRepository()
	userID := account.(*xpaseto.PasetoPayload).AccountID
	accountEntity, err := accountRepo.GetAccountByID(ctx, userID)
	if err != nil {
		log.Errorw("Failed to get account by ID", zap.Error(err))
		return nil, err
	}
	return &out.GetProfileResponse{
		Email:     accountEntity.Email,
		CreatedAt: accountEntity.CreatedAt.Format(time.RFC3339),
		UpdatedAt: accountEntity.UpdatedAt.Format(time.RFC3339),
	}, nil
}
