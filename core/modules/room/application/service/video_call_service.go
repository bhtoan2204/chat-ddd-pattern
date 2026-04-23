package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	appCtx "wechat-clone/core/context"
	apptypes "wechat-clone/core/modules/room/application/types"
	"wechat-clone/core/modules/room/constant"
	"wechat-clone/core/modules/room/domain/entity"
	roomrepos "wechat-clone/core/modules/room/domain/repos"
	sharedcache "wechat-clone/core/shared/infra/cache"
	"wechat-clone/core/shared/infra/lock"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type VideoCallService interface {
	EnsureRoomMember(ctx context.Context, roomID, actorID string) error
	GetActiveCall(ctx context.Context, query apptypes.GetActiveVideoCallQuery) (*apptypes.VideoCallSessionResult, error)
	StartCall(ctx context.Context, command apptypes.StartVideoCallCommand) (*apptypes.VideoCallSessionResult, error)
	JoinCall(ctx context.Context, command apptypes.JoinVideoCallCommand) (*apptypes.VideoCallSessionResult, error)
	LeaveCall(ctx context.Context, command apptypes.LeaveVideoCallCommand) (*apptypes.VideoCallSessionResult, error)
	EndCall(ctx context.Context, command apptypes.EndVideoCallCommand) (*apptypes.VideoCallSessionResult, error)
	RelaySignal(ctx context.Context, command apptypes.RelayVideoCallSignalCommand) (*apptypes.VideoCallSignalResult, error)
}

var ErrVideoCallActiveSessionAlreadyExists = errors.New("video call active session already exists")

type videoCallSessionStore interface {
	Get(ctx context.Context, roomID string) (*entity.VideoCallSession, bool, error)
	Save(ctx context.Context, session *entity.VideoCallSession) error
	Delete(ctx context.Context, roomID string) error
}

type videoCallService struct {
	baseRepo roomrepos.Repos
	locker   lock.Lock
	store    videoCallSessionStore
}

func NewVideoCallService(appContext *appCtx.AppContext, baseRepo roomrepos.Repos) VideoCallService {
	if appContext == nil || baseRepo == nil {
		return nil
	}

	return &videoCallService{
		baseRepo: baseRepo,
		locker:   appContext.Locker(),
		store:    newVideoCallSessionCacheStore(appContext.GetCache()),
	}
}

func (s *videoCallService) EnsureRoomMember(ctx context.Context, roomID, actorID string) error {
	_, err := s.requireRoomMember(ctx, roomID, actorID)
	return stackErr.Error(err)
}

func (s *videoCallService) GetActiveCall(ctx context.Context, query apptypes.GetActiveVideoCallQuery) (*apptypes.VideoCallSessionResult, error) {
	if _, err := s.requireRoomMember(ctx, query.RoomID, query.ActorID); err != nil {
		return nil, stackErr.Error(err)
	}

	session, found, err := s.store.Get(ctx, query.RoomID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if !found || session == nil || !session.IsActive() {
		return nil, nil
	}
	return buildVideoCallSessionResult(session), nil
}

func (s *videoCallService) StartCall(ctx context.Context, command apptypes.StartVideoCallCommand) (*apptypes.VideoCallSessionResult, error) {
	return withVideoCallRoomLock(ctx, s.locker, command.RoomID, func() (*apptypes.VideoCallSessionResult, error) {
		if _, err := s.requireRoomMember(ctx, command.RoomID, command.ActorID); err != nil {
			return nil, stackErr.Error(err)
		}

		existing, found, err := s.store.Get(ctx, command.RoomID)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		if found && existing != nil && existing.IsActive() {
			return nil, stackErr.Error(ErrVideoCallActiveSessionAlreadyExists)
		}

		session, err := entity.NewVideoCallSession(uuid.NewString(), strings.TrimSpace(command.RoomID), strings.TrimSpace(command.ActorID), time.Now().UTC())
		if err != nil {
			return nil, stackErr.Error(err)
		}
		if err := s.store.Save(ctx, session); err != nil {
			return nil, stackErr.Error(err)
		}
		return buildVideoCallSessionResult(session), nil
	})
}

func (s *videoCallService) JoinCall(ctx context.Context, command apptypes.JoinVideoCallCommand) (*apptypes.VideoCallSessionResult, error) {
	return withVideoCallRoomLock(ctx, s.locker, command.RoomID, func() (*apptypes.VideoCallSessionResult, error) {
		if _, err := s.requireRoomMember(ctx, command.RoomID, command.ActorID); err != nil {
			return nil, stackErr.Error(err)
		}

		session, err := s.requireActiveSession(ctx, command.RoomID, command.SessionID)
		if err != nil {
			return nil, stackErr.Error(err)
		}

		if err := session.Join(command.ActorID, time.Now().UTC()); err != nil {
			return nil, stackErr.Error(err)
		}
		if err := s.store.Save(ctx, session); err != nil {
			return nil, stackErr.Error(err)
		}
		return buildVideoCallSessionResult(session), nil
	})
}

func (s *videoCallService) LeaveCall(ctx context.Context, command apptypes.LeaveVideoCallCommand) (*apptypes.VideoCallSessionResult, error) {
	return withVideoCallRoomLock(ctx, s.locker, command.RoomID, func() (*apptypes.VideoCallSessionResult, error) {
		if _, err := s.requireRoomMember(ctx, command.RoomID, command.ActorID); err != nil {
			return nil, stackErr.Error(err)
		}

		session, err := s.requireActiveSession(ctx, command.RoomID, command.SessionID)
		if err != nil {
			return nil, stackErr.Error(err)
		}

		ended, err := session.Leave(command.ActorID, time.Now().UTC())
		if err != nil {
			return nil, stackErr.Error(err)
		}
		if ended {
			if err := s.store.Delete(ctx, command.RoomID); err != nil {
				return nil, stackErr.Error(err)
			}
		} else if err := s.store.Save(ctx, session); err != nil {
			return nil, stackErr.Error(err)
		}
		return buildVideoCallSessionResult(session), nil
	})
}

func (s *videoCallService) EndCall(ctx context.Context, command apptypes.EndVideoCallCommand) (*apptypes.VideoCallSessionResult, error) {
	return withVideoCallRoomLock(ctx, s.locker, command.RoomID, func() (*apptypes.VideoCallSessionResult, error) {
		if _, err := s.requireRoomMember(ctx, command.RoomID, command.ActorID); err != nil {
			return nil, stackErr.Error(err)
		}

		session, err := s.requireActiveSession(ctx, command.RoomID, command.SessionID)
		if err != nil {
			return nil, stackErr.Error(err)
		}
		if !session.HasParticipant(command.ActorID) && strings.TrimSpace(session.StartedByAccountID) != strings.TrimSpace(command.ActorID) {
			return nil, stackErr.Error(entity.ErrVideoCallParticipantNotFound)
		}

		if err := session.End(command.ActorID, time.Now().UTC()); err != nil {
			return nil, stackErr.Error(err)
		}
		if err := s.store.Delete(ctx, command.RoomID); err != nil {
			return nil, stackErr.Error(err)
		}
		return buildVideoCallSessionResult(session), nil
	})
}

func (s *videoCallService) RelaySignal(ctx context.Context, command apptypes.RelayVideoCallSignalCommand) (*apptypes.VideoCallSignalResult, error) {
	if _, err := s.requireRoomMember(ctx, command.RoomID, command.ActorID); err != nil {
		return nil, stackErr.Error(err)
	}
	if err := entity.ValidateVideoCallSignalType(command.SignalType); err != nil {
		return nil, stackErr.Error(err)
	}
	if err := entity.ValidateVideoCallSignalTarget(command.TargetAccountID); err != nil {
		return nil, stackErr.Error(err)
	}

	session, err := s.requireActiveSession(ctx, command.RoomID, command.SessionID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if !session.HasParticipant(command.ActorID) {
		return nil, stackErr.Error(entity.ErrVideoCallParticipantNotFound)
	}

	return &apptypes.VideoCallSignalResult{
		SessionID:       session.SessionID,
		RoomID:          session.RoomID,
		SenderAccountID: strings.TrimSpace(command.ActorID),
		TargetAccountID: strings.TrimSpace(command.TargetAccountID),
		SignalType:      strings.TrimSpace(command.SignalType),
		SignalPayload:   cloneRawJSON(command.SignalPayloadRaw),
	}, nil
}

func (s *videoCallService) requireRoomMember(ctx context.Context, roomID, actorID string) (*entity.RoomMemberEntity, error) {
	roomAgg, err := s.baseRepo.RoomAggregateRepository().Load(ctx, strings.TrimSpace(roomID))
	if err != nil {
		return nil, stackErr.Error(err)
	}

	for _, member := range roomAgg.Members() {
		if member == nil {
			continue
		}
		if strings.TrimSpace(member.AccountID) == strings.TrimSpace(actorID) {
			return member, nil
		}
	}
	return nil, stackErr.Error(entity.ErrRoomMemberRequired)
}

func (s *videoCallService) requireActiveSession(ctx context.Context, roomID, sessionID string) (*entity.VideoCallSession, error) {
	session, found, err := s.store.Get(ctx, roomID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if !found || session == nil || !session.IsActive() {
		return nil, stackErr.Error(entity.ErrVideoCallAlreadyEnded)
	}
	if strings.TrimSpace(sessionID) != "" && strings.TrimSpace(session.SessionID) != strings.TrimSpace(sessionID) {
		return nil, stackErr.Error(fmt.Errorf("video call session mismatch: %s", strings.TrimSpace(sessionID)))
	}
	return session, nil
}

func buildVideoCallSessionResult(session *entity.VideoCallSession) *apptypes.VideoCallSessionResult {
	if session == nil {
		return nil
	}

	result := &apptypes.VideoCallSessionResult{
		SessionID:             session.SessionID,
		RoomID:                session.RoomID,
		Status:                session.Status,
		StartedByAccountID:    session.StartedByAccountID,
		ParticipantAccountIDs: append([]string(nil), session.ParticipantAccountIDs...),
		StartedAt:             session.StartedAt.UTC().Format(time.RFC3339),
		UpdatedAt:             session.UpdatedAt.UTC().Format(time.RFC3339),
		EndedByAccountID:      session.EndedByAccountID,
	}
	if session.EndedAt != nil {
		result.EndedAt = session.EndedAt.UTC().Format(time.RFC3339)
	}
	return result
}

func withVideoCallRoomLock[T any](ctx context.Context, locker lock.Lock, roomID string, fn func() (T, error)) (T, error) {
	return lock.WithLocks(ctx, locker, []string{strings.TrimSpace(roomID)}, constant.DefaultVideoCallLockOptions(), fn)
}

func cloneRawJSON(payload json.RawMessage) json.RawMessage {
	if len(payload) == 0 {
		return nil
	}
	cloned := make([]byte, len(payload))
	copy(cloned, payload)
	return json.RawMessage(cloned)
}

type videoCallSessionCacheStore struct {
	cache sharedcache.Cache
}

func newVideoCallSessionCacheStore(cache sharedcache.Cache) videoCallSessionStore {
	if cache == nil {
		return nil
	}
	return &videoCallSessionCacheStore{cache: cache}
}

func (s *videoCallSessionCacheStore) Get(ctx context.Context, roomID string) (*entity.VideoCallSession, bool, error) {
	if s == nil || s.cache == nil {
		return nil, false, nil
	}

	data, err := s.cache.Get(ctx, videoCallSessionCacheKey(roomID))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, stackErr.Error(err)
	}

	session := &entity.VideoCallSession{}
	if err := json.Unmarshal(data, session); err != nil {
		return nil, false, stackErr.Error(fmt.Errorf("unmarshal video call session cache: %w", err))
	}
	return session, true, nil
}

func (s *videoCallSessionCacheStore) Save(ctx context.Context, session *entity.VideoCallSession) error {
	if s == nil || s.cache == nil || session == nil {
		return nil
	}
	return stackErr.Error(s.cache.SetObject(ctx, videoCallSessionCacheKey(session.RoomID), session, constant.VideoCallSessionTTL))
}

func (s *videoCallSessionCacheStore) Delete(ctx context.Context, roomID string) error {
	if s == nil || s.cache == nil {
		return nil
	}
	return stackErr.Error(s.cache.Delete(ctx, videoCallSessionCacheKey(roomID)))
}

func videoCallSessionCacheKey(roomID string) string {
	return fmt.Sprintf("room:video_call:%s", strings.TrimSpace(roomID))
}
