package service

import (
	"context"
	"strings"
	"time"

	"go-socket/core/modules/room/application/projection"
	roomsupport "go-socket/core/modules/room/application/support"
	apptypes "go-socket/core/modules/room/application/types"
	"go-socket/core/shared/pkg/stackErr"
)

type MessageQueryService interface {
	ListMessages(ctx context.Context, accountID string, query apptypes.ListMessagesQuery) ([]apptypes.MessageResult, error)
}

type messageQueryService struct {
	readRepos projection.QueryRepos
}

func newMessageQueryService(readRepos projection.QueryRepos) MessageQueryService {
	return &messageQueryService{readRepos: readRepos}
}

func (s *messageQueryService) ListMessages(ctx context.Context, accountID string, query apptypes.ListMessagesQuery) ([]apptypes.MessageResult, error) {
	limit := query.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var beforeAt *time.Time
	if strings.TrimSpace(query.BeforeAt) != "" {
		if parsed, err := time.Parse(time.RFC3339, query.BeforeAt); err == nil {
			beforeAt = &parsed
		}
	}

	messages, err := s.readRepos.MessageReadRepository().ListMessages(ctx, accountID, query.RoomID, projection.MessageListOptions{
		Limit:     limit,
		BeforeID:  query.BeforeID,
		BeforeAt:  beforeAt,
		Ascending: query.Ascending,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	out := make([]apptypes.MessageResult, 0, len(messages))
	for _, message := range messages {
		item, err := roomsupport.BuildMessageResult(ctx, s.readRepos, accountID, message)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		out = append(out, *item)
	}
	return out, nil
}
