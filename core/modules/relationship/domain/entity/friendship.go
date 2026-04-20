package entity

import (
	"fmt"
	"time"
)

type Friendship struct {
	ID                   string
	UserLowID            string
	UserHighID           string
	CreatedAt            time.Time
	CreatedFromRequestID *string
}

func NewFriendship(
	id string,
	userA string,
	userB string,
	createdFromRequestID *string,
	now time.Time,
) (*Friendship, error) {
	if id == "" {
		return nil, fmt.Errorf("friendship id is required")
	}
	if userA == "" || userB == "" {
		return nil, fmt.Errorf("friendship users are required")
	}
	if userA == userB {
		return nil, fmt.Errorf("cannot create friendship with self")
	}

	userLowID, userHighID := normalizePair(userA, userB)

	return &Friendship{
		ID:                   id,
		UserLowID:            userLowID,
		UserHighID:           userHighID,
		CreatedAt:            now,
		CreatedFromRequestID: createdFromRequestID,
	}, nil
}

func (f *Friendship) HasUser(userID string) bool {
	return f.UserLowID == userID || f.UserHighID == userID
}

func (f *Friendship) OtherUserID(userID string) (string, error) {
	switch userID {
	case f.UserLowID:
		return f.UserHighID, nil
	case f.UserHighID:
		return f.UserLowID, nil
	default:
		return "", fmt.Errorf("user does not belong to this friendship")
	}
}
