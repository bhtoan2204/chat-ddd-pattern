package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"wechat-clone/core/modules/notification/domain/entity"
	notificationrepos "wechat-clone/core/modules/notification/domain/repos"
	"wechat-clone/core/shared/pkg/stackErr"
	sharedwebpush "wechat-clone/core/shared/pkg/webpush"
)

//go:generate mockgen -package=service -destination=push_delivery_service_mock.go -source=push_delivery_service.go
type PushDeliveryService interface {
	SendNotification(ctx context.Context, notification *entity.NotificationEntity) error
}

type pushDeliveryService struct {
	pushSubscriptions notificationrepos.PushSubscriptionRepository
	webPush           sharedwebpush.WebPushService
}

func NewPushDeliveryService(
	pushSubscriptions notificationrepos.PushSubscriptionRepository,
	webPush sharedwebpush.WebPushService,
) PushDeliveryService {
	if pushSubscriptions == nil || webPush == nil {
		return nil
	}

	return &pushDeliveryService{
		pushSubscriptions: pushSubscriptions,
		webPush:           webPush,
	}
}

func (s *pushDeliveryService) SendNotification(ctx context.Context, notification *entity.NotificationEntity) error {
	accountID := strings.TrimSpace(notification.AccountID)
	subscriptions, err := s.pushSubscriptions.ListPushSubscriptionsByAccountID(ctx, accountID)
	if err != nil {
		return stackErr.Error(fmt.Errorf("list push subscriptions failed: %w", err))
	}
	if len(subscriptions) == 0 {
		return nil
	}

	payload, err := buildWebPushPayload(notification)
	if err != nil {
		return stackErr.Error(err)
	}

	items := make([]sharedwebpush.Subscription, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		item, mapErr := mapSubscription(subscription)
		if mapErr != nil {
			return stackErr.Error(mapErr)
		}
		items = append(items, item)
	}

	if err := s.webPush.SendMany(ctx, payload, items); err != nil {
		return stackErr.Error(fmt.Errorf("send webpush failed: %w", err))
	}
	return nil
}

type webPushPayload struct {
	Title string                 `json:"title"`
	Body  string                 `json:"body"`
	Data  map[string]interface{} `json:"data,omitempty"`
}

func buildWebPushPayload(notification *entity.NotificationEntity) ([]byte, error) {
	if notification == nil {
		return nil, stackErr.Error(fmt.Errorf("notification is nil"))
	}

	title := strings.TrimSpace(notification.Subject)
	if title == "" {
		title = "PayChat"
	}

	body := strings.TrimSpace(notification.Body)
	if body == "" {
		body = "Ban co thong bao moi."
	}

	payload := webPushPayload{
		Title: title,
		Body:  body,
		Data: map[string]interface{}{
			"notification_id": notification.ID,
			"account_id":      notification.AccountID,
			"url":             resolveNotificationURL(notification),
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, stackErr.Error(fmt.Errorf("marshal webpush payload failed: %w", err))
	}
	return data, nil
}

func mapSubscription(subscription *entity.PushSubscription) (sharedwebpush.Subscription, error) {
	if subscription == nil {
		return sharedwebpush.Subscription{}, stackErr.Error(fmt.Errorf("push subscription is nil"))
	}

	item := sharedwebpush.Subscription{
		Endpoint: strings.TrimSpace(subscription.Endpoint),
	}

	if err := json.Unmarshal([]byte(subscription.Keys), &item.Keys); err != nil {
		return sharedwebpush.Subscription{}, stackErr.Error(fmt.Errorf("unmarshal push subscription keys failed: %w", err))
	}
	if err := item.Validate(); err != nil {
		return sharedwebpush.Subscription{}, stackErr.Error(err)
	}

	return item, nil
}

func resolveNotificationURL(notification *entity.NotificationEntity) string {
	if notification == nil {
		return "/"
	}

	roomID := strings.TrimSpace(notification.RoomID)
	if roomID != "" {
		return "/chat/" + roomID
	}

	switch {
	case strings.TrimSpace(notification.LastMessageID) != "":
		return "/notifications/" + strings.TrimSpace(notification.ID)
	default:
		return "/notifications"
	}
}
