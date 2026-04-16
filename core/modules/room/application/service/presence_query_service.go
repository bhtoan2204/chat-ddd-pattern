package service

import (
	"context"
	"strings"

	appCtx "go-socket/core/context"
	apptypes "go-socket/core/modules/room/application/types"
	"go-socket/core/shared/pkg/stackErr"

	"github.com/redis/go-redis/v9"
)

type PresenceQueryService interface {
	GetPresence(ctx context.Context, query apptypes.GetPresenceQuery) (*apptypes.PresenceResult, error)
}

type presenceQueryService struct {
	redis *redis.Client
}

func newPresenceQueryService(appCtx *appCtx.AppContext) PresenceQueryService {
	return &presenceQueryService{redis: appCtx.GetRedisClient()}
}

func (s *presenceQueryService) GetPresence(ctx context.Context, query apptypes.GetPresenceQuery) (*apptypes.PresenceResult, error) {
	accountID := query.AccountID
	if s.redis == nil {
		return &apptypes.PresenceResult{AccountID: accountID, Status: "offline"}, nil
	}

	exists, err := s.redis.Exists(ctx, chatPresenceKey(accountID)).Result()
	if err != nil {
		return nil, stackErr.Error(err)
	}

	status := "offline"
	if exists > 0 {
		status = "online"
	}
	return &apptypes.PresenceResult{AccountID: accountID, Status: status}, nil
}

func chatPresenceKey(accountID string) string {
	return "chat:presence:" + strings.TrimSpace(accountID)
}
