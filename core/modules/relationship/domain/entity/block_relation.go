package entity

import (
	"fmt"
	"time"
)

type BlockRelation struct {
	ID        string
	BlockerID string
	BlockedID string
	Reason    *string
	CreatedAt time.Time
}

func NewBlockRelation(id string, blockerID string, blockedID string, reason *string, now time.Time) (*BlockRelation, error) {
	if id == "" {
		return nil, fmt.Errorf("block relation id is required")
	}
	if blockerID == "" {
		return nil, fmt.Errorf("blocker id is required")
	}
	if blockedID == "" {
		return nil, fmt.Errorf("blocked id is required")
	}
	if blockerID == blockedID {
		return nil, fmt.Errorf("cannot block self")
	}

	return &BlockRelation{
		ID:        id,
		BlockerID: blockerID,
		BlockedID: blockedID,
		Reason:    reason,
		CreatedAt: now,
	}, nil
}
