package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"wechat-clone/core/modules/notification/application/support"
	"wechat-clone/core/modules/notification/domain/aggregate"
	notificationrepos "wechat-clone/core/modules/notification/domain/repos"
	notificationtypes "wechat-clone/core/modules/notification/types"
	sharedevents "wechat-clone/core/shared/contracts/events"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

type PaymentNotificationService interface {
	NotifyWithdrawalRequested(ctx context.Context, payload sharedevents.PaymentWithdrawalRequestedEvent) error
	NotifyWithdrawalSucceeded(ctx context.Context, payload sharedevents.PaymentSucceededEvent) error
	NotifyWithdrawalFailed(ctx context.Context, payload sharedevents.PaymentFailedEvent) error
}

type paymentNotificationService struct {
	baseRepo notificationrepos.Repos
	realtime RealtimeService
	push     PushDeliveryService
}

type generalNotificationSpec struct {
	NotificationID string
	AccountID      string
	Type           notificationtypes.NotificationType
	Subject        string
	Body           string
	OccurredAt     time.Time
}

func newPaymentNotificationService(
	baseRepo notificationrepos.Repos,
	realtime RealtimeService,
	push PushDeliveryService,
) PaymentNotificationService {
	if baseRepo == nil {
		return nil
	}

	return &paymentNotificationService{
		baseRepo: baseRepo,
		realtime: realtime,
		push:     push,
	}
}

func (s *paymentNotificationService) NotifyWithdrawalRequested(ctx context.Context, payload sharedevents.PaymentWithdrawalRequestedEvent) error {
	paymentID := firstNonEmpty(payload.PaymentID, payload.TransactionID)
	return stackErr.Error(s.createGeneralNotificationAndEmit(ctx, generalNotificationSpec{
		NotificationID: aggregate.PaymentNotificationID(notificationtypes.NotificationTypeWithdrawalRequested, paymentID, payload.DebitAccountID),
		AccountID:      strings.TrimSpace(payload.DebitAccountID),
		Type:           notificationtypes.NotificationTypeWithdrawalRequested,
		Subject:        "Withdrawal requested",
		Body:           fmt.Sprintf("Your withdrawal request for %d %s has been received and is being processed.", payload.Amount, payload.Currency),
		OccurredAt:     payload.RequestedAt,
	}))
}

func (s *paymentNotificationService) NotifyWithdrawalSucceeded(ctx context.Context, payload sharedevents.PaymentSucceededEvent) error {
	if strings.TrimSpace(payload.Workflow) != "WITHDRAWAL" {
		return nil
	}
	paymentID := firstNonEmpty(payload.PaymentID, payload.TransactionID)
	return stackErr.Error(s.createGeneralNotificationAndEmit(ctx, generalNotificationSpec{
		NotificationID: aggregate.PaymentNotificationID(notificationtypes.NotificationTypeWithdrawalSucceeded, paymentID, payload.DebitAccountID),
		AccountID:      strings.TrimSpace(payload.DebitAccountID),
		Type:           notificationtypes.NotificationTypeWithdrawalSucceeded,
		Subject:        "Withdrawal completed",
		Body:           fmt.Sprintf("Your withdrawal of %d %s completed successfully.", payload.Amount, payload.Currency),
		OccurredAt:     payload.SucceededAt,
	}))
}

func (s *paymentNotificationService) NotifyWithdrawalFailed(ctx context.Context, payload sharedevents.PaymentFailedEvent) error {
	if strings.TrimSpace(payload.Workflow) != "WITHDRAWAL" {
		return nil
	}
	paymentID := firstNonEmpty(payload.PaymentID, payload.TransactionID)
	return stackErr.Error(s.createGeneralNotificationAndEmit(ctx, generalNotificationSpec{
		NotificationID: aggregate.PaymentNotificationID(notificationtypes.NotificationTypeWithdrawalFailed, paymentID, payload.DebitAccountID),
		AccountID:      strings.TrimSpace(payload.DebitAccountID),
		Type:           notificationtypes.NotificationTypeWithdrawalFailed,
		Subject:        "Withdrawal failed",
		Body:           fmt.Sprintf("Your withdrawal of %d %s failed and the reserved balance was returned.", payload.Amount, payload.Currency),
		OccurredAt:     payload.OccurredAt,
	}))
}

func (s *paymentNotificationService) createGeneralNotificationAndEmit(ctx context.Context, spec generalNotificationSpec) error {
	if strings.TrimSpace(spec.NotificationID) == "" || strings.TrimSpace(spec.AccountID) == "" {
		return nil
	}

	notificationRepo := s.baseRepo.NotificationRepository()
	if _, err := notificationRepo.Load(ctx, spec.NotificationID); err == nil {
		return nil
	} else if !errors.Is(err, notificationrepos.ErrNotificationNotFound) {
		return stackErr.Error(err)
	}

	notificationAgg, err := aggregate.NewNotificationAggregate(spec.NotificationID)
	if err != nil {
		return stackErr.Error(err)
	}

	if err := notificationAgg.Create(
		spec.AccountID,
		spec.Type,
		spec.Subject,
		spec.Body,
		spec.OccurredAt,
	); err != nil {
		return stackErr.Error(err)
	}
	if err := notificationRepo.Save(ctx, notificationAgg); err != nil {
		return stackErr.Error(fmt.Errorf("save payment notification failed: %w", err))
	}

	snapshot, err := notificationAgg.Snapshot()
	if err != nil {
		return stackErr.Error(err)
	}
	unreadCount, err := notificationRepo.CountUnread(ctx, spec.AccountID)
	if err != nil {
		return stackErr.Error(err)
	}
	if s.realtime != nil {
		if emitErr := s.realtime.EmitMessage(ctx, support.NewRealtimeNotificationPayload(notificationtypes.RealtimeEventNotificationUpsert, snapshot, unreadCount)); emitErr != nil {
			logging.FromContext(ctx).Warnw("emit payment notification realtime failed", zap.Error(emitErr))
		}
	}
	if s.push != nil {
		if pushErr := s.push.SendNotification(ctx, snapshot); pushErr != nil {
			logging.FromContext(ctx).Warnw("send payment notification webpush failed", zap.Error(pushErr))
		}
	}

	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
