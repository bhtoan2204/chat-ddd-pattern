package service

import (
	"context"

	"go-socket/core/modules/room/application/projection"
	apptypes "go-socket/core/modules/room/application/types"

	"github.com/redis/go-redis/v9"
)

type ChatQueryService struct {
	conversations *conversationQueryService
	messages      *messageQueryService
	mentions      *mentionQueryService
	presence      *presenceQueryService
}

func NewChatQueryService(readRepos projection.QueryRepos, redis *redis.Client) *ChatQueryService {
	return &ChatQueryService{
		conversations: newConversationQueryService(readRepos),
		messages:      newMessageQueryService(readRepos),
		mentions:      newMentionQueryService(readRepos),
		presence:      newPresenceQueryService(redis),
	}
}

func (s *ChatQueryService) ListConversations(ctx context.Context, accountID string, query apptypes.ListConversationsQuery) ([]apptypes.ConversationResult, error) {
	return s.conversations.ListConversations(ctx, accountID, query)
}

func (s *ChatQueryService) GetConversation(ctx context.Context, accountID string, query apptypes.GetConversationQuery) (*apptypes.ConversationResult, error) {
	return s.conversations.GetConversation(ctx, accountID, query)
}

func (s *ChatQueryService) ListMessages(ctx context.Context, accountID string, query apptypes.ListMessagesQuery) ([]apptypes.MessageResult, error) {
	return s.messages.ListMessages(ctx, accountID, query)
}

func (s *ChatQueryService) SearchMentionCandidates(ctx context.Context, accountID string, query apptypes.SearchMentionCandidatesQuery) ([]apptypes.MentionCandidateResult, error) {
	return s.mentions.SearchMentionCandidates(ctx, accountID, query)
}

func (s *ChatQueryService) GetPresence(ctx context.Context, query apptypes.GetPresenceQuery) (*apptypes.PresenceResult, error) {
	return s.presence.GetPresence(ctx, query)
}
