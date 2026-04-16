package service

import (
	"context"

	appCtx "go-socket/core/context"
	"go-socket/core/modules/room/constant"
	"go-socket/core/modules/room/types"
	"go-socket/core/shared/pkg/pubsub"
	"go-socket/core/shared/pkg/stackErr"
)

type RealtimeService interface {
	EmitMessage(ctx context.Context, message types.MessagePayload) error
}

type realtimeService struct {
	localPublisher *pubsub.Bus
}

func newRealtimeService(appCtx *appCtx.AppContext) RealtimeService {
	return &realtimeService{
		localPublisher: appCtx.LocalBus(),
	}
}

func (s *realtimeService) EmitMessage(ctx context.Context, message types.MessagePayload) error {
	if err := s.localPublisher.Publish(ctx, constant.RealtimeMessageTopic, message); err != nil {
		return stackErr.Error(err)
	}
	return nil
}
