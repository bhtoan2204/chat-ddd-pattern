package models

import "time"

type FollowRelation struct {
	ID         string    `gorm:"column:id;type:varchar(36);primaryKey"`
	FollowerID string    `gorm:"column:follower_id;type:varchar(36);not null;uniqueIndex:uq_follows_pair,priority:1;index:idx_follows_follower_created,priority:1"`
	FolloweeID string    `gorm:"column:followee_id;type:varchar(36);not null;uniqueIndex:uq_follows_pair,priority:2;index:idx_follows_followee_created,priority:1"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime;index:idx_follows_follower_created,priority:2,sort:desc;index:idx_follows_followee_created,priority:2,sort:desc"`
}

func (FollowRelation) TableName() string {
	return "relationship_follows"
}
