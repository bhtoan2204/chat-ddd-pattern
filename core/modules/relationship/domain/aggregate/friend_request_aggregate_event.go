package aggregate

import (
	"time"
)

type EventFriendRequestCreated struct {
	RequesterID string
	AddresseeID string
	Message     string
	CreatedAt   time.Time
}

type EventFriendRequestAccept struct {
	AcceptedAt time.Time
}

type EventFriendRequestReject struct {
	Reason     *string
	RejectedAt time.Time
}

type EventFriendRequestCancel struct {
	CancelAt time.Time
}
