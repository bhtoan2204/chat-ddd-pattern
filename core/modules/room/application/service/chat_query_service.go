package service

import (
	"context"

	"go-socket/core/modules/room/application/projection"
	apptypes "go-socket/core/modules/room/application/types"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -package=service -destination=chat_query_service_mock.go -source=chat_query_service.go
type ChatQueryService interface {
	ConversationQueryService
	MessageQueryService
	MentionQueryService
	PresenceQueryService
}

type chatQueryService struct {
	conversations ConversationQueryService
	messages      MessageQueryService
	mentions      MentionQueryService
	presence      PresenceQueryService
}

func NewChatQueryService(readRepos projection.QueryRepos, redis *redis.Client) ChatQueryService {
	return &chatQueryService{
		conversations: newConversationQueryService(readRepos),
		messages:      newMessageQueryService(readRepos),
		mentions:      newMentionQueryService(readRepos),
		presence:      newPresenceQueryService(redis),
	}
}

func (s *chatQueryService) ListConversations(ctx context.Context, accountID string, query apptypes.ListConversationsQuery) ([]apptypes.ConversationResult, error) {
	return s.conversations.ListConversations(ctx, accountID, query)
}

func (s *chatQueryService) GetConversation(ctx context.Context, accountID string, query apptypes.GetConversationQuery) (*apptypes.ConversationResult, error) {
	return s.conversations.GetConversation(ctx, accountID, query)
}

func (s *chatQueryService) ListMessages(ctx context.Context, accountID string, query apptypes.ListMessagesQuery) ([]apptypes.MessageResult, error) {
	return s.messages.ListMessages(ctx, accountID, query)
}

func (s *chatQueryService) SearchMentionCandidates(ctx context.Context, accountID string, query apptypes.SearchMentionCandidatesQuery) ([]apptypes.MentionCandidateResult, error) {
	return s.mentions.SearchMentionCandidates(ctx, accountID, query)
}

func (s *chatQueryService) GetPresence(ctx context.Context, query apptypes.GetPresenceQuery) (*apptypes.PresenceResult, error) {
	return s.presence.GetPresence(ctx, query)
}
