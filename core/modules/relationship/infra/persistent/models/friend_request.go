package models

import "time"

type FriendRequestStatus string

const (
	FriendRequestStatusPending   FriendRequestStatus = "PENDING"
	FriendRequestStatusAccepted  FriendRequestStatus = "ACCEPTED"
	FriendRequestStatusRejected  FriendRequestStatus = "REJECTED"
	FriendRequestStatusCancelled FriendRequestStatus = "CANCELLED"
	FriendRequestStatusExpired   FriendRequestStatus = "EXPIRED"
)

type FriendRequest struct {
	ID             string              `gorm:"column:id;type:varchar(36);primaryKey"`
	RequesterID    string              `gorm:"column:requester_id;type:varchar(36);not null;index:idx_friend_requests_requester_status_created,priority:1"`
	AddresseeID    string              `gorm:"column:addressee_id;type:varchar(36);not null;index:idx_friend_requests_addressee_status_created,priority:1"`
	Status         FriendRequestStatus `gorm:"column:status;type:varchar(20);not null;index:idx_friend_requests_requester_status_created,priority:2;index:idx_friend_requests_addressee_status_created,priority:2"`
	Message        *string             `gorm:"column:message;type:varchar(500)"`
	CreatedAt      time.Time           `gorm:"column:created_at;type:timestamptz;not null;autoCreateTime;index:idx_friend_requests_requester_status_created,priority:3,sort:desc;index:idx_friend_requests_addressee_status_created,priority:3,sort:desc"`
	RespondedAt    *time.Time          `gorm:"column:responded_at;type:timestamptz"`
	ExpiredAt      *time.Time          `gorm:"column:expired_at;type:timestamptz"`
	CancelledAt    *time.Time          `gorm:"column:cancelled_at;type:timestamptz"`
	RejectedReason *string             `gorm:"column:rejected_reason;type:varchar(255)"`
}

func (FriendRequest) TableName() string {
	return "relationship_friend_requests"
}
