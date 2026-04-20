package models

import "time"

type UserRelationshipCounters struct {
	UserID          string    `gorm:"column:user_id;type:varchar(36);primaryKey"`
	FriendsCount    int64     `gorm:"column:friends_count;not null;default:0"`
	FollowersCount  int64     `gorm:"column:followers_count;not null;default:0"`
	FollowingCount  int64     `gorm:"column:following_count;not null;default:0"`
	BlockedCount    int64     `gorm:"column:blocked_count;not null;default:0"`
	PendingInCount  int64     `gorm:"column:pending_in_count;not null;default:0"`
	PendingOutCount int64     `gorm:"column:pending_out_count;not null;default:0"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamptz;not null;autoUpdateTime"`
}

func (UserRelationshipCounters) TableName() string {
	return "relationship_user_relationship_counters"
}
