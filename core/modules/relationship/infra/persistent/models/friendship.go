package models

import "time"

type Friendship struct {
	ID                   string    `gorm:"column:id;type:varchar(36);primaryKey"`
	UserLowID            string    `gorm:"column:user_low_id;type:varchar(36);not null;uniqueIndex:uq_friendships_pair,priority:1;index:idx_friendships_user_low_created,priority:1"`
	UserHighID           string    `gorm:"column:user_high_id;type:varchar(36);not null;uniqueIndex:uq_friendships_pair,priority:2;index:idx_friendships_user_high_created,priority:1"`
	CreatedAt            time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime;index:idx_friendships_user_low_created,priority:2,sort:desc;index:idx_friendships_user_high_created,priority:2,sort:desc"`
	CreatedFromRequestID *string   `gorm:"column:created_from_request_id;type:varchar(36)"`
}

func (Friendship) TableName() string {
	return "relationship_friendships"
}
