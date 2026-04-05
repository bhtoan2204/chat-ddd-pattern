package command

import (
	"context"
	"fmt"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	accountservice "go-socket/core/modules/account/application/service"
	"go-socket/core/modules/account/domain/entity"
	repos "go-socket/core/modules/account/domain/repos"
	valueobject "go-socket/core/modules/account/domain/value_object"
	"go-socket/core/shared/infra/xpaseto"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/hasher"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type registerHandler struct {
	baseRepo         repos.Repos
	hasher           hasher.Hasher
	paseto           xpaseto.PasetoService
	aggregateService *accountservice.AccountAggregateService
}

func NewRegisterHandler(appCtx *appCtx.AppContext, baseRepo repos.Repos, aggregateService *accountservice.AccountAggregateService) cqrs.Handler[*in.RegisterRequest, *out.RegisterResponse] {
	return &registerHandler{
		baseRepo:         baseRepo,
		hasher:           appCtx.GetHasher(),
		paseto:           appCtx.GetPaseto(),
		aggregateService: aggregateService,
	}
}

func (u *registerHandler) Handle(ctx context.Context, req *in.RegisterRequest) (*out.RegisterResponse, error) {
	log := logging.FromContext(ctx).Named("Register")
	accountRepo := u.baseRepo.AccountRepository()

	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		log.Errorw("Failed to create email", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	exists, err := accountRepo.IsEmailExists(ctx, email.Value())
	if err != nil {
		log.Errorw("Failed to check existing account", zap.Error(err))
		return nil, stackErr.Error(ErrCheckAccountFailed)
	}
	if exists {
		log.Errorw("Account already exists", zap.String("email", email.Value()))
		return nil, stackErr.Error(ErrAccountExists)
	}

	password, err := valueobject.NewPassword(req.Password)
	if err != nil {
		log.Errorw("Failed to create password", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	hashedPassword, err := u.hasher.Hash(ctx, password.Value())
	if err != nil {
		log.Errorw("Failed to hash password", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	hashedPasswordVO, err := valueobject.NewPassword(hashedPassword)
	if err != nil {
		log.Errorw("Failed to create hashed password value object", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	now := time.Now().UTC()
	newAccountEntity, err := entity.NewAccount(uuid.New().String(), email, hashedPasswordVO, req.DisplayName, now)
	if err != nil {
		log.Errorw("Failed to create account entity", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	if txErr := u.baseRepo.WithTransaction(ctx, func(txRepos repos.Repos) error {
		if err := txRepos.AccountRepository().CreateAccount(ctx, newAccountEntity); err != nil {
			log.Errorw("Failed to create account", zap.Error(err))
			return stackErr.Error(fmt.Errorf("create account failed: %v", err))
		}

		if err := u.aggregateService.PublishAccountCreated(ctx, txRepos.AccountOutboxEventsRepository(), newAccountEntity); err != nil {
			return fmt.Errorf("publish account created event failed: %v", err)
		}
		return nil
	}); txErr != nil {
		log.Errorw("Failed to register account", zap.Error(txErr))
		return nil, stackErr.Error(txErr)
	}

	token, expiresAt, err := u.paseto.GenerateToken(ctx, newAccountEntity)
	if err != nil {
		log.Errorw("Failed to generate token", zap.Error(err))
		return nil, stackErr.Error(fmt.Errorf("generate token failed: %v", err))
	}

	return &out.RegisterResponse{
		Token:     token,
		ExpiresAt: expiresAt.UnixMilli(),
	}, nil
}
