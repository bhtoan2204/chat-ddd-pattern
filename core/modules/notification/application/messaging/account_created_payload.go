package messaging

import (
	"encoding/json"
	"fmt"
	stackerr "go-socket/core/shared/pkg/stackErr"
	"strings"
)

type accountCreatedPayload struct {
	AccountID string `json:"AccountID"`
	Email     string `json:"Email"`
	CreatedAt string `json:"CreatedAt"`
}

func parseAccountCreatedPayload(raw []byte) (*accountCreatedPayload, error) {
	if len(raw) == 0 {
		return nil, stackerr.Error(fmt.Errorf("event_data is empty"))
	}

	var payload accountCreatedPayload
	if err := json.Unmarshal(raw, &payload); err == nil {
		return &payload, nil
	}

	var encoded string
	if err := json.Unmarshal(raw, &encoded); err != nil {
		return nil, stackerr.Error(fmt.Errorf("unmarshal event_data as string failed: %w", err))
	}

	encoded = strings.TrimSpace(encoded)
	if encoded == "" {
		return nil, stackerr.Error(fmt.Errorf("event_data is empty"))
	}

	if err := json.Unmarshal([]byte(encoded), &payload); err != nil {
		return nil, stackerr.Error(fmt.Errorf("unmarshal encoded event_data failed: %w", err))
	}

	return &payload, nil
}
