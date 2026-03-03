package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"go-socket/config"
	"go-socket/core/modules/notification/application/adapter"
	infraMessaging "go-socket/core/shared/infra/messaging"
	stackerr "go-socket/core/shared/pkg/stackErr"
	"strings"
)

type MessageHandler interface {
	Start() error
	Stop() error
}

type messageHandler struct {
	consumer    []infraMessaging.Consumer
	emailSender adapter.EmailSender
}

const accountCreatedEventName = "account.created"

type accountOutboxMessage struct {
	ID            int64           `json:"ID"`
	AggregateID   string          `json:"AGGREGATE_ID"`
	AggregateType string          `json:"AGGREGATE_TYPE"`
	Version       int64           `json:"VERSION"`
	EventName     string          `json:"EVENT_NAME"`
	EventData     json.RawMessage `json:"EVENT_DATA"`
	CreatedAt     string          `json:"CREATED_AT"`
}

func NewMessageHandler(cfg *config.Config, emailSender adapter.EmailSender) (MessageHandler, error) {
	if emailSender == nil {
		return nil, stackerr.Error(fmt.Errorf("email sender can not be nil"))
	}

	instance := &messageHandler{
		emailSender: emailSender,
		consumer:    make([]infraMessaging.Consumer, 0),
	}

	consumeTopics := []string{cfg.KafkaConfig.KafkaNotificationConsumer.AccountTopic}
	mapHandler := map[string]infraMessaging.Handler{
		fmt.Sprintf("notification-%s-handler", strings.ToLower(cfg.KafkaConfig.KafkaNotificationConsumer.AccountTopic)): func(ctx context.Context, value []byte) error {
			return instance.handleAccountCreated(ctx, value)
		},
	}

	for _, topic := range consumeTopics {
		consumer, err := infraMessaging.NewConsumer(&infraMessaging.Config{
			Servers:      cfg.KafkaConfig.KafkaServers,
			Group:        cfg.KafkaConfig.KafkaNotificationConsumer.NotificationGroup,
			OffsetReset:  cfg.KafkaConfig.KafkaOffsetReset,
			ConsumeTopic: []string{topic},
			HandlerName:  fmt.Sprintf("notification-%s-handler", strings.ToLower(topic)),
			DLQ:          true,
		})
		if err != nil {
			return nil, stackerr.Error(err)
		}
		consumer.SetHandler(mapHandler[fmt.Sprintf("notification-%s-handler", strings.ToLower(topic))])
		instance.consumer = append(instance.consumer, consumer)
	}

	return instance, nil
}

func (h *messageHandler) Start() error {
	for _, consumer := range h.consumer {
		consumer.Read(h.processMessage(consumer))
	}
	return nil
}

func (h *messageHandler) Stop() error {
	for _, consumer := range h.consumer {
		consumer.Stop()
	}
	return nil
}
