package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/payment/application/dto/in"
	paymentservice "wechat-clone/core/modules/payment/application/service"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/contracts"
	sharedevents "wechat-clone/core/shared/contracts/events"
	infraMessaging "wechat-clone/core/shared/infra/messaging"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

//go:generate mockgen -package=messaging -destination=message_handler_mock.go -source=message_handler.go
type MessageHandler interface {
	Start() error
	Stop() error
}

type messageHandler struct {
	consumer       []infraMessaging.Consumer
	paymentService paymentservice.PaymentCommandService
}

func NewMessageHandler(
	cfg *config.Config,
	_ *appCtx.AppContext,
	paymentService paymentservice.PaymentCommandService,
) (MessageHandler, error) {
	instance := &messageHandler{
		consumer:       make([]infraMessaging.Consumer, 0, 1),
		paymentService: paymentService,
	}

	topic := strings.TrimSpace(cfg.KafkaConfig.KafkaPaymentConsumer.LedgerOutboxTopic)
	if topic == "" {
		return instance, nil
	}

	consumer, err := infraMessaging.NewConsumer(&infraMessaging.Config{
		Servers:      cfg.KafkaConfig.KafkaServers,
		Group:        cfg.KafkaConfig.KafkaPaymentConsumer.PaymentGroup,
		OffsetReset:  cfg.KafkaConfig.KafkaOffsetReset,
		ConsumeTopic: []string{topic},
		HandlerName:  fmt.Sprintf("payment-%s-handler", strings.ToLower(topic)),
		DLQ:          true,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	consumer.SetHandler(func(ctx context.Context, value []byte) error {
		return instance.handleLedgerOutboxEvent(ctx, value)
	})
	instance.consumer = append(instance.consumer, consumer)

	return instance, nil
}

func (h *messageHandler) Start() error {
	for _, consumer := range h.consumer {
		consumer.Read(infraMessaging.WrapConsumerCallback(consumer, "Handle payment message failed"))
	}
	return nil
}

func (h *messageHandler) Stop() error {
	infraMessaging.StopConsumers(h.consumer)
	return nil
}

func (h *messageHandler) handleLedgerOutboxEvent(ctx context.Context, value []byte) error {
	log := logging.FromContext(ctx).Named("PaymentLedgerEvent")

	var event contracts.OutboxMessage
	if err := json.Unmarshal(value, &event); err != nil {
		return stackErr.Error(fmt.Errorf("unmarshal ledger outbox event failed: %w", err))
	}

	log.Infow("handle ledger outbox event",
		zap.String("event_name", event.EventName),
		zap.String("aggregate_id", event.AggregateID),
	)

	switch event.EventName {
	case sharedevents.EventLedgerPaymentReconciliationFailed:
		payload, err := unmarshalLedgerPaymentReconciliationFailedPayload(event.EventData)
		if err != nil {
			return stackErr.Error(err)
		}
		return h.refundFailedLedgerReconciliation(ctx, payload)
	default:
		return nil
	}
}

func unmarshalLedgerPaymentReconciliationFailedPayload(data json.RawMessage) (sharedevents.LedgerPaymentReconciliationFailedEvent, error) {
	var payload sharedevents.LedgerPaymentReconciliationFailedEvent
	if err := contracts.UnmarshalEventData(data, &payload); err != nil {
		return sharedevents.LedgerPaymentReconciliationFailedEvent{}, stackErr.Error(fmt.Errorf("unmarshal ledger payment reconciliation failed payload failed: %w", err))
	}
	return payload, nil
}

func (h *messageHandler) refundFailedLedgerReconciliation(ctx context.Context, payload sharedevents.LedgerPaymentReconciliationFailedEvent) error {
	if h.paymentService == nil {
		return stackErr.Error(fmt.Errorf("payment command service is required"))
	}
	transactionID := strings.TrimSpace(payload.TransactionID)
	if transactionID == "" {
		transactionID = strings.TrimSpace(payload.PaymentID)
	}
	if transactionID == "" {
		return stackErr.Error(fmt.Errorf("ledger reconciliation failure missing transaction_id"))
	}

	_, err := h.paymentService.RefundPayment(ctx, &in.RefundPaymentRequest{
		Provider:      strings.TrimSpace(payload.Provider),
		TransactionID: transactionID,
		Reason:        refundReason(payload),
	})
	if err != nil {
		return stackErr.Error(err)
	}
	return nil
}

func refundReason(payload sharedevents.LedgerPaymentReconciliationFailedEvent) string {
	reason := strings.TrimSpace(payload.Reason)
	if reason == "" {
		reason = "ledger reconciliation failed"
	}
	return reason
}
