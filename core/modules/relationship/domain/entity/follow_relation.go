package entity

import (
	"fmt"
	"time"
)

type FollowRelation struct {
	ID         string
	FollowerID string
	FolloweeID string
	CreatedAt  time.Time
}

func NewFollowRelation(id string, followerID string, followeeID string, now time.Time) (*FollowRelation, error) {
	if id == "" {
		return nil, fmt.Errorf("follow relation id is required")
	}
	if followerID == "" {
		return nil, fmt.Errorf("follower id is required")
	}
	if followeeID == "" {
		return nil, fmt.Errorf("followee id is required")
	}
	if followerID == followeeID {
		return nil, fmt.Errorf("cannot follow self")
	}

	return &FollowRelation{
		ID:         id,
		FollowerID: followerID,
		FolloweeID: followeeID,
		CreatedAt:  now,
	}, nil
}
