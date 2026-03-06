package types

type NotificationType string

const (
	NotificationTypeAccountCreated NotificationType = "account.created"
)

func (t NotificationType) String() string {
	return string(t)
}
