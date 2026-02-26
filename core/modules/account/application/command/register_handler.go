package command

import (
	"context"
	"errors"
	"fmt"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/domain/entity"
	repos "go-socket/core/modules/account/domain/repos"
	valueobject "go-socket/core/modules/account/domain/value_object"
	"go-socket/core/shared/contracts/events"
	"go-socket/core/shared/infra/xpaseto"
	eventpkg "go-socket/core/shared/pkg/event"
	"go-socket/core/shared/pkg/hasher"
	"go-socket/core/shared/pkg/logging"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type registerHandler struct {
	baseRepo repos.Repos
	hasher   hasher.Hasher
	paseto   xpaseto.PasetoService
}

func NewRegisterHandler(appCtx *appCtx.AppContext, baseRepo repos.Repos) RegisterHandler {
	return &registerHandler{
		baseRepo: baseRepo,
		hasher:   appCtx.GetHasher(),
		paseto:   appCtx.GetPaseto(),
	}
}

func (u *registerHandler) Handle(ctx context.Context, req *in.RegisterRequest) (*out.RegisterResponse, error) {
	log := logging.FromContext(ctx).Named("Register")
	accountRepo := u.baseRepo.AccountRepository()
	_, err := accountRepo.GetAccountByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorw("Failed to check existing account", zap.Error(err))
		return nil, ErrCheckAccountFailed
	}
	if err == nil {
		log.Errorw("Account already exists", zap.String("email", req.Email))
		return nil, ErrAccountExists
	}

	password, err := valueobject.NewPassword(req.Password)
	if err != nil {
		log.Errorw("Failed to create password", zap.Error(err))
		return nil, err
	}

	hashedPassword, err := u.hasher.Hash(ctx, password.Value())
	if err != nil {
		log.Errorw("Failed to hash password", zap.Error(err))
		return nil, err
	}

	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		log.Errorw("Failed to create email", zap.Error(err))
		return nil, err
	}

	hashedPasswordVO, err := valueobject.NewPassword(hashedPassword)
	if err != nil {
		log.Errorw("Failed to create hashed password value object", zap.Error(err))
		return nil, err
	}

	newAccountEntity := &entity.Account{
		ID:       uuid.New().String(),
		Email:    email,
		Password: hashedPasswordVO,
	}

	if txErr := u.baseRepo.WithTransaction(ctx, func(txRepos repos.Repos) error {
		if err := txRepos.AccountRepository().CreateAccount(ctx, newAccountEntity); err != nil {
			log.Errorw("Failed to create account", zap.Error(err))
			return fmt.Errorf("create account failed: %w", err)
		}

		payload := &events.AccountCreatedEvent{
			AccountID: newAccountEntity.ID,
			Email:     newAccountEntity.Email.Value(),
			CreatedAt: time.Now(),
		}
		evt := eventpkg.Event{
			AggregateID:   newAccountEntity.ID,
			AggregateType: "account",
			Version:       1,
			EventName:     payload.GetName(),
			EventData:     payload,
			CreatedAt:     payload.CreatedAt.Unix(),
		}
		publisher := eventpkg.NewPublisher(txRepos.AccountOutboxEventsRepository())
		if err := publisher.Publish(ctx, evt); err != nil {
			log.Errorw("Failed to append account created event", zap.Error(err))
			return fmt.Errorf("append account created event failed: %w", err)
		}
		return nil
	}); txErr != nil {
		log.Errorw("Failed to register account", zap.Error(txErr))
		return nil, txErr
	}

	token, expiresAt, err := u.paseto.GenerateToken(ctx, newAccountEntity)
	if err != nil {
		log.Errorw("Failed to generate token", zap.Error(err))
		return nil, fmt.Errorf("generate token failed: %w", err)
	}

	return &out.RegisterResponse{
		Token:     token,
		ExpiresAt: expiresAt.UnixMilli(),
	}, nil
}
