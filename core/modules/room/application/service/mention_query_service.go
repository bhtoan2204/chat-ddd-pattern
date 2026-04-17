package service

import (
	"context"
	"errors"
	"strings"

	"wechat-clone/core/modules/room/application/projection"
	apptypes "wechat-clone/core/modules/room/application/types"
	"wechat-clone/core/modules/room/infra/projection/cassandra/views"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/samber/lo"
)

type MentionQueryService interface {
	SearchMentionCandidates(ctx context.Context, accountID string, query apptypes.SearchMentionCandidatesQuery) ([]apptypes.MentionCandidateResult, error)
}

type mentionQueryService struct {
	readRepos projection.QueryRepos
}

func newMentionQueryService(readRepos projection.QueryRepos) MentionQueryService {
	return &mentionQueryService{readRepos: readRepos}
}

func (s *mentionQueryService) SearchMentionCandidates(ctx context.Context, accountID string, query apptypes.SearchMentionCandidatesQuery) ([]apptypes.MentionCandidateResult, error) {
	roomID := strings.TrimSpace(query.RoomID)
	if roomID == "" {
		return nil, stackErr.Error(errors.New("room_id is required"))
	}

	room, err := s.readRepos.RoomReadRepository().GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if room == nil || !strings.EqualFold(strings.TrimSpace(room.RoomType), "group") {
		return nil, stackErr.Error(errors.New("mentions are supported only in group rooms"))
	}

	member, err := s.readRepos.RoomMemberReadRepository().GetRoomMemberByAccount(ctx, roomID, accountID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if member == nil {
		return nil, stackErr.Error(errors.New("viewer is not a member of this room"))
	}

	candidates, err := s.readRepos.RoomMemberReadRepository().SearchMentionCandidates(ctx, projection.MentionCandidateSearch{
		RoomID:           roomID,
		Keyword:          query.Query,
		ExcludeAccountID: accountID,
		Limit:            query.Limit,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	results := lo.FilterMap(candidates, func(candidate *views.MentionCandidateView, _ int) (apptypes.MentionCandidateResult, bool) {
		if candidate == nil {
			return apptypes.MentionCandidateResult{}, false
		}

		return apptypes.MentionCandidateResult{
			AccountID:       candidate.AccountID,
			DisplayName:     resolveMentionCandidateDisplayName(candidate),
			Username:        strings.TrimSpace(candidate.Username),
			AvatarObjectKey: strings.TrimSpace(candidate.AvatarObjectKey),
		}, true
	})
	if len(results) == 0 {
		return nil, nil
	}
	return results, nil
}

func resolveMentionCandidateDisplayName(candidate *views.MentionCandidateView) string {
	if candidate == nil {
		return ""
	}

	switch {
	case strings.TrimSpace(candidate.DisplayName) != "":
		return strings.TrimSpace(candidate.DisplayName)
	case strings.TrimSpace(candidate.Username) != "":
		return strings.TrimSpace(candidate.Username)
	default:
		return strings.TrimSpace(candidate.AccountID)
	}
}
