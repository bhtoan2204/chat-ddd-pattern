package command

import "errors"

var (
	ErrAccountNotFound            = errors.New("account not found")
	ErrSavePushSubscriptionFailed = errors.New("save push subscription failed")
)
