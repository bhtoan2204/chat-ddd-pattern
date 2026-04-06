package query

import (
	"context"
	"time"

	appCtx "go-socket/core/context"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/dto/out"
	"go-socket/core/modules/account/application/service"
	repos "go-socket/core/modules/account/domain/repos"
	"go-socket/core/shared/infra/storage"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

const avatarPresignedURLTTL = 15 * time.Minute

type getAvatarHandler struct {
	accountRepo repos.AccountRepository
	storage     storage.Storage
}

func NewGetAvatarHandler(appCtx *appCtx.AppContext, baseRepo repos.Repos, services service.Services) cqrs.Handler[*in.GetAvatarRequest, *out.GetAvatarResponse] {
	return &getAvatarHandler{
		accountRepo: baseRepo.AccountRepository(),
		storage:     appCtx.GetStorage(),
	}
}

func (u *getAvatarHandler) Handle(ctx context.Context, req *in.GetAvatarRequest) (*out.GetAvatarResponse, error) {
	log := logging.FromContext(ctx).Named("GetAvatar")

	accountEntity, err := u.accountRepo.GetAccountByID(ctx, req.AccountID)
	if err != nil {
		log.Errorw("Failed to get account by ID", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	if accountEntity.AvatarObjectKey == nil || *accountEntity.AvatarObjectKey == "" {
		return &out.GetAvatarResponse{}, nil
	}

	url, err := u.storage.PresignedGetObjectURL(ctx, *accountEntity.AvatarObjectKey, avatarPresignedURLTTL)
	if err != nil {
		log.Errorw("Failed to generate presigned avatar URL", zap.Error(err))
		return nil, stackErr.Error(err)
	}

	return &out.GetAvatarResponse{
		URL:       url,
		ExpiresAt: time.Now().UTC().Add(avatarPresignedURLTTL).Format(time.RFC3339),
	}, nil
}
