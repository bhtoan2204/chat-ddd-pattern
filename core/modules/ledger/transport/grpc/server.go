package grpc

import (
	"context"
	"fmt"
	"strings"

	ledgerapp "wechat-clone/core/modules/ledger/application/service"
	"wechat-clone/core/shared/contracts"
	sharedevents "wechat-clone/core/shared/contracts/events"
	sharedlock "wechat-clone/core/shared/infra/lock"
	"wechat-clone/core/shared/pkg/stackErr"
	ledgerv1 "wechat-clone/core/shared/transport/grpc/gen/ledger/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ledgerPaymentGRPCServer struct {
	ledgerv1.LedgerPaymentServiceServer
	locker       sharedlock.Lock
	service      ledgerapp.PaymentEventService
	feeAccountID string
}

func NewPaymentServer(
	locker sharedlock.Lock,
	service ledgerapp.PaymentEventService,
	feeAccountID string,
) ledgerv1.LedgerPaymentServiceServer {
	return &ledgerPaymentGRPCServer{
		locker:       locker,
		service:      service,
		feeAccountID: strings.TrimSpace(feeAccountID),
	}
}

func (s *ledgerPaymentGRPCServer) ApplyPaymentEvent(ctx context.Context, req *ledgerv1.ApplyPaymentEventRequest) (*ledgerv1.ApplyPaymentEventResponse, error) {
	applied, err := s.applyPaymentEvent(ctx, strings.TrimSpace(req.GetEventName()), []byte(req.GetEventDataJson()))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	return &ledgerv1.ApplyPaymentEventResponse{Applied: applied}, nil
}

func (s *ledgerPaymentGRPCServer) applyPaymentEvent(ctx context.Context, eventName string, raw []byte) (bool, error) {
	switch eventName {
	case sharedevents.EventPaymentCreated, sharedevents.EventPaymentCheckoutSessionCreated:
		return false, nil
	case sharedevents.EventPaymentWithdrawalRequested:
		var payload sharedevents.PaymentWithdrawalRequestedEvent
		if err := contracts.UnmarshalEventData(raw, &payload); err != nil {
			return false, stackErr.Error(fmt.Errorf("decode payment withdrawal requested payload failed: %w", err))
		}
		return true, stackErr.Error(s.withLedgerLocks(ctx, withdrawalRequestedLockKeys(payload, s.feeAccountID), func() error {
			return stackErr.Error(s.service.HandleWithdrawalRequested(ctx, payload))
		}))
	case sharedevents.EventPaymentSucceeded:
		var payload sharedevents.PaymentSucceededEvent
		if err := contracts.UnmarshalEventData(raw, &payload); err != nil {
			return false, stackErr.Error(fmt.Errorf("decode payment succeeded payload failed: %w", err))
		}
		if strings.EqualFold(strings.TrimSpace(payload.Workflow), "WITHDRAWAL") {
			return false, nil
		}
		return true, stackErr.Error(s.withLedgerLocks(ctx, succeededLockKeys(payload, s.feeAccountID), func() error {
			return stackErr.Error(s.service.HandleSucceeded(ctx, payload))
		}))
	case sharedevents.EventPaymentFailed:
		var payload sharedevents.PaymentFailedEvent
		if err := contracts.UnmarshalEventData(raw, &payload); err != nil {
			return false, stackErr.Error(fmt.Errorf("decode payment failed payload failed: %w", err))
		}
		if !strings.EqualFold(strings.TrimSpace(payload.Workflow), "WITHDRAWAL") {
			return false, nil
		}
		return true, stackErr.Error(s.withLedgerLocks(ctx, failedLockKeys(payload, s.feeAccountID), func() error {
			return stackErr.Error(s.service.HandleFailed(ctx, payload))
		}))
	case sharedevents.EventPaymentRefunded:
		var payload sharedevents.PaymentRefundedEvent
		if err := contracts.UnmarshalEventData(raw, &payload); err != nil {
			return false, stackErr.Error(fmt.Errorf("decode payment refunded payload failed: %w", err))
		}
		return true, stackErr.Error(s.withLedgerLocks(ctx, reversedLockKeys(payload.ClearingAccountKey, payload.CreditAccountID, payload.FeeAmount, s.feeAccountID), func() error {
			return stackErr.Error(s.service.HandleRefunded(ctx, payload))
		}))
	case sharedevents.EventPaymentChargeback:
		var payload sharedevents.PaymentChargebackEvent
		if err := contracts.UnmarshalEventData(raw, &payload); err != nil {
			return false, stackErr.Error(fmt.Errorf("decode payment chargeback payload failed: %w", err))
		}
		return true, stackErr.Error(s.withLedgerLocks(ctx, reversedLockKeys(payload.ClearingAccountKey, payload.CreditAccountID, payload.FeeAmount, s.feeAccountID), func() error {
			return stackErr.Error(s.service.HandleChargeback(ctx, payload))
		}))
	default:
		return false, nil
	}
}

func (s *ledgerPaymentGRPCServer) withLedgerLocks(ctx context.Context, keys []string, fn func() error) error {
	if s.service == nil {
		return nil
	}

	opts := sharedlock.DefaultMultiLockOptions()
	opts.KeyPrefix = ledgerapp.LedgerAccountLockKeyPrefix

	_, err := sharedlock.WithLocks(ctx, s.locker, keys, opts, func() (struct{}, error) {
		return struct{}{}, fn()
	})
	if err != nil {
		return stackErr.Error(err)
	}
	return nil
}

func withdrawalRequestedLockKeys(payload sharedevents.PaymentWithdrawalRequestedEvent, feeAccountID string) []string {
	keys := []string{
		strings.TrimSpace(payload.DebitAccountID),
		ledgerClearingAccountID(firstNonEmpty(payload.ClearingAccountKey, providerClearingAccountKey(payload.Provider))),
	}
	if payload.FeeAmount > 0 {
		keys = append(keys, strings.TrimSpace(feeAccountID))
	}
	return normalizeLockKeys(keys)
}

func succeededLockKeys(payload sharedevents.PaymentSucceededEvent, feeAccountID string) []string {
	keys := []string{
		ledgerClearingAccountID(payload.ClearingAccountKey),
		strings.TrimSpace(payload.CreditAccountID),
	}
	if payload.FeeAmount > 0 {
		keys = append(keys, strings.TrimSpace(feeAccountID))
	}
	return normalizeLockKeys(keys)
}

func failedLockKeys(payload sharedevents.PaymentFailedEvent, feeAccountID string) []string {
	keys := []string{
		strings.TrimSpace(payload.DebitAccountID),
		ledgerClearingAccountID(payload.ClearingAccountKey),
	}
	if payload.FeeAmount > 0 {
		keys = append(keys, strings.TrimSpace(feeAccountID))
	}
	return normalizeLockKeys(keys)
}

func reversedLockKeys(clearingAccountKey, creditAccountID string, feeAmount int64, feeAccountID string) []string {
	keys := []string{
		ledgerClearingAccountID(clearingAccountKey),
		strings.TrimSpace(creditAccountID),
	}
	if feeAmount > 0 {
		keys = append(keys, strings.TrimSpace(feeAccountID))
	}
	return normalizeLockKeys(keys)
}

func normalizeLockKeys(keys []string) []string {
	items := make([]string, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		items = append(items, key)
	}
	return items
}

func ledgerClearingAccountID(clearingAccountKey string) string {
	clearingAccountKey = strings.ToLower(strings.TrimSpace(clearingAccountKey))
	if clearingAccountKey == "" {
		return ""
	}
	return fmt.Sprintf("ledger:clearing:%s", clearingAccountKey)
}

func providerClearingAccountKey(provider string) string {
	provider = strings.ToLower(strings.TrimSpace(provider))
	if provider == "" {
		return ""
	}
	return fmt.Sprintf("provider:%s", provider)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func mapGRPCError(err error) error {
	return status.Error(codes.Internal, err.Error())
}
