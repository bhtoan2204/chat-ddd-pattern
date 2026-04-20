package models

import "time"

type RelationOutboxEvent struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`
	AggregateID   string    `gorm:"column:aggregate_id;type:varchar(1024);not null"`
	AggregateType string    `gorm:"column:aggregate_type;type:varchar(255);not null;index:idx_relationship_outbox_aggregate,priority:1"`
	Version       int64     `gorm:"column:version;not null;index:idx_relationship_outbox_aggregate,priority:3"`
	EventName     string    `gorm:"column:event_name;type:varchar(255);not null"`
	EventData     string    `gorm:"column:event_data;type:text;not null"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime;index:idx_relationship_outbox_created_at;index:idx_relationship_outbox_aggregate,priority:2"`
}

func (RelationOutboxEvent) TableName() string {
	return "relationship_outbox_events"
}
