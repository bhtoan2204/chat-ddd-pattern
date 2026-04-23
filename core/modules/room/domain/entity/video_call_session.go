package entity

import (
	"errors"
	"strings"
	"time"

	"wechat-clone/core/shared/pkg/stackErr"
)

const (
	VideoCallStatusActive = "active"
	VideoCallStatusEnded  = "ended"
)

var (
	ErrVideoCallSessionIDRequired     = errors.New("video call session_id is required")
	ErrVideoCallRoomRequired          = errors.New("video call room_id is required")
	ErrVideoCallActorRequired         = errors.New("video call actor_id is required")
	ErrVideoCallAlreadyEnded          = errors.New("video call is already ended")
	ErrVideoCallParticipantNotFound   = errors.New("video call participant is not in session")
	ErrVideoCallSignalTypeRequired    = errors.New("video call signal_type is required")
	ErrVideoCallTargetAccountRequired = errors.New("video call target_account_id is required")
)

type VideoCallSession struct {
	SessionID             string     `json:"session_id"`
	RoomID                string     `json:"room_id"`
	Status                string     `json:"status"`
	StartedByAccountID    string     `json:"started_by_account_id"`
	ParticipantAccountIDs []string   `json:"participant_account_ids"`
	StartedAt             time.Time  `json:"started_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	EndedAt               *time.Time `json:"ended_at,omitempty"`
	EndedByAccountID      string     `json:"ended_by_account_id,omitempty"`
}

func NewVideoCallSession(sessionID, roomID, startedByAccountID string, now time.Time) (*VideoCallSession, error) {
	sessionID = strings.TrimSpace(sessionID)
	roomID = strings.TrimSpace(roomID)
	startedByAccountID = strings.TrimSpace(startedByAccountID)

	switch {
	case sessionID == "":
		return nil, stackErr.Error(ErrVideoCallSessionIDRequired)
	case roomID == "":
		return nil, stackErr.Error(ErrVideoCallRoomRequired)
	case startedByAccountID == "":
		return nil, stackErr.Error(ErrVideoCallActorRequired)
	}

	now = normalizeRoomTime(now)
	return &VideoCallSession{
		SessionID:             sessionID,
		RoomID:                roomID,
		Status:                VideoCallStatusActive,
		StartedByAccountID:    startedByAccountID,
		ParticipantAccountIDs: []string{startedByAccountID},
		StartedAt:             now,
		UpdatedAt:             now,
	}, nil
}

func (s *VideoCallSession) IsActive() bool {
	return s != nil && s.Status == VideoCallStatusActive && s.EndedAt == nil
}

func (s *VideoCallSession) HasParticipant(accountID string) bool {
	if s == nil {
		return false
	}
	accountID = strings.TrimSpace(accountID)
	if accountID == "" {
		return false
	}
	for _, participantID := range s.ParticipantAccountIDs {
		if participantID == accountID {
			return true
		}
	}
	return false
}

func (s *VideoCallSession) Join(accountID string, now time.Time) error {
	if s == nil {
		return stackErr.Error(ErrVideoCallSessionIDRequired)
	}
	if !s.IsActive() {
		return stackErr.Error(ErrVideoCallAlreadyEnded)
	}
	accountID = strings.TrimSpace(accountID)
	if accountID == "" {
		return stackErr.Error(ErrVideoCallActorRequired)
	}
	if s.HasParticipant(accountID) {
		s.UpdatedAt = normalizeRoomTime(now)
		return nil
	}

	s.ParticipantAccountIDs = append(s.ParticipantAccountIDs, accountID)
	s.UpdatedAt = normalizeRoomTime(now)
	return nil
}

func (s *VideoCallSession) Leave(accountID string, now time.Time) (bool, error) {
	if s == nil {
		return false, stackErr.Error(ErrVideoCallSessionIDRequired)
	}
	if !s.IsActive() {
		return false, stackErr.Error(ErrVideoCallAlreadyEnded)
	}
	accountID = strings.TrimSpace(accountID)
	if accountID == "" {
		return false, stackErr.Error(ErrVideoCallActorRequired)
	}
	if !s.HasParticipant(accountID) {
		return false, stackErr.Error(ErrVideoCallParticipantNotFound)
	}

	filtered := make([]string, 0, len(s.ParticipantAccountIDs))
	for _, participantID := range s.ParticipantAccountIDs {
		if participantID == accountID {
			continue
		}
		filtered = append(filtered, participantID)
	}
	s.ParticipantAccountIDs = filtered
	s.UpdatedAt = normalizeRoomTime(now)

	if len(s.ParticipantAccountIDs) == 0 {
		if err := s.End(accountID, now); err != nil {
			return false, stackErr.Error(err)
		}
		return true, nil
	}

	return false, nil
}

func (s *VideoCallSession) End(actorID string, now time.Time) error {
	if s == nil {
		return stackErr.Error(ErrVideoCallSessionIDRequired)
	}
	if !s.IsActive() {
		return stackErr.Error(ErrVideoCallAlreadyEnded)
	}
	actorID = strings.TrimSpace(actorID)
	if actorID == "" {
		return stackErr.Error(ErrVideoCallActorRequired)
	}

	endedAt := normalizeRoomTime(now)
	s.Status = VideoCallStatusEnded
	s.UpdatedAt = endedAt
	s.EndedAt = &endedAt
	s.EndedByAccountID = actorID
	return nil
}

func ValidateVideoCallSignalTarget(targetAccountID string) error {
	if strings.TrimSpace(targetAccountID) == "" {
		return stackErr.Error(ErrVideoCallTargetAccountRequired)
	}
	return nil
}

func ValidateVideoCallSignalType(signalType string) error {
	if strings.TrimSpace(signalType) == "" {
		return stackErr.Error(ErrVideoCallSignalTypeRequired)
	}
	return nil
}
