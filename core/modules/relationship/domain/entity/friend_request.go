package entity

import (
	"fmt"
	"time"
	"wechat-clone/core/shared/pkg/stackErr"
)

type FriendRequestStatus string

func (r *FriendRequestStatus) String() string {
	return string(*r)
}

const (
	FriendRequestStatusPending   FriendRequestStatus = "PENDING"
	FriendRequestStatusAccepted  FriendRequestStatus = "ACCEPTED"
	FriendRequestStatusRejected  FriendRequestStatus = "REJECTED"
	FriendRequestStatusCancelled FriendRequestStatus = "CANCELLED"
	FriendRequestStatusExpired   FriendRequestStatus = "EXPIRED"
)

type FriendRequest struct {
	ID             string
	RequesterID    string
	AddresseeID    string
	Status         FriendRequestStatus
	Message        *string
	CreatedAt      time.Time
	RespondedAt    *time.Time
	ExpiredAt      *time.Time
	CancelledAt    *time.Time
	RejectedReason *string
}

func NewFriendRequest(
	id string,
	requesterID string,
	addresseeID string,
	message *string,
	now time.Time,
) (*FriendRequest, error) {
	if id == "" {
		return nil, fmt.Errorf("friend request id is required")
	}
	if requesterID == "" {
		return nil, fmt.Errorf("requester id is required")
	}
	if addresseeID == "" {
		return nil, fmt.Errorf("addressee id is required")
	}
	if requesterID == addresseeID {
		return nil, fmt.Errorf("cannot send friend request to self")
	}

	return &FriendRequest{
		ID:          id,
		RequesterID: requesterID,
		AddresseeID: addresseeID,
		Status:      FriendRequestStatusPending,
		Message:     message,
		CreatedAt:   now,
	}, nil
}

func (r *FriendRequest) Accept(now time.Time) error {
	if r.Status != FriendRequestStatusPending {
		return stackErr.Error(fmt.Errorf("friend request is not pending"))
	}

	r.Status = FriendRequestStatusAccepted
	r.RespondedAt = &now

	return nil
}

func (r *FriendRequest) Reject(reason *string, now time.Time) error {
	if r.Status != FriendRequestStatusPending {
		return fmt.Errorf("friend request is not pending")
	}

	r.Status = FriendRequestStatusRejected
	r.RespondedAt = &now
	r.RejectedReason = reason

	return nil
}

func (r *FriendRequest) Cancel(now time.Time) error {
	if r.Status != FriendRequestStatusPending {
		return fmt.Errorf("friend request is not pending")
	}

	r.Status = FriendRequestStatusCancelled
	r.CancelledAt = &now

	return nil
}

func (r *FriendRequest) Expire(now time.Time) error {
	if r.Status != FriendRequestStatusPending {
		return fmt.Errorf("friend request is not pending")
	}

	r.Status = FriendRequestStatusExpired
	r.ExpiredAt = &now

	return nil
}

func (r *FriendRequest) IsPending() bool {
	return r.Status == FriendRequestStatusPending
}
