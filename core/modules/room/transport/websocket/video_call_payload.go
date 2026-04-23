package socket

import "encoding/json"

type videoCallSessionRequest struct {
	SessionID string `json:"session_id"`
}

type videoCallSignalRequest struct {
	SessionID       string          `json:"session_id"`
	TargetAccountID string          `json:"target_account_id"`
	SignalType      string          `json:"signal_type"`
	SignalPayload   json.RawMessage `json:"signal_payload"`
}
