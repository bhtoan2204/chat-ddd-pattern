package models

import "time"

type BlockRelation struct {
	ID        string    `gorm:"column:id;type:varchar(36);primaryKey"`
	BlockerID string    `gorm:"column:blocker_id;type:varchar(36);not null;uniqueIndex:uq_blocks_pair,priority:1;index:idx_blocks_blocker_created,priority:1"`
	BlockedID string    `gorm:"column:blocked_id;type:varchar(36);not null;uniqueIndex:uq_blocks_pair,priority:2;index:idx_blocks_blocked_created,priority:1"`
	Reason    *string   `gorm:"column:reason;type:varchar(255)"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime;index:idx_blocks_blocker_created,priority:2,sort:desc;index:idx_blocks_blocked_created,priority:2,sort:desc"`
}

func (BlockRelation) TableName() string {
	return "relationship_blocks"
}
