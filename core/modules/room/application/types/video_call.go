package types

import "encoding/json"

type StartVideoCallCommand struct {
	RoomID  string
	ActorID string
}

type JoinVideoCallCommand struct {
	RoomID    string
	SessionID string
	ActorID   string
}

type LeaveVideoCallCommand struct {
	RoomID    string
	SessionID string
	ActorID   string
}

type EndVideoCallCommand struct {
	RoomID    string
	SessionID string
	ActorID   string
}

type GetActiveVideoCallQuery struct {
	RoomID  string
	ActorID string
}

type RelayVideoCallSignalCommand struct {
	RoomID           string
	SessionID        string
	ActorID          string
	TargetAccountID  string
	SignalType       string
	SignalPayloadRaw json.RawMessage
}

type VideoCallSessionResult struct {
	SessionID             string   `json:"session_id"`
	RoomID                string   `json:"room_id"`
	Status                string   `json:"status"`
	StartedByAccountID    string   `json:"started_by_account_id"`
	ParticipantAccountIDs []string `json:"participant_account_ids"`
	StartedAt             string   `json:"started_at"`
	UpdatedAt             string   `json:"updated_at"`
	EndedAt               string   `json:"ended_at,omitempty"`
	EndedByAccountID      string   `json:"ended_by_account_id,omitempty"`
}

type VideoCallSignalResult struct {
	SessionID       string          `json:"session_id"`
	RoomID          string          `json:"room_id"`
	SenderAccountID string          `json:"sender_account_id"`
	TargetAccountID string          `json:"target_account_id"`
	SignalType      string          `json:"signal_type"`
	SignalPayload   json.RawMessage `json:"signal_payload,omitempty"`
}
