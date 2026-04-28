package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	accountprojection "wechat-clone/core/modules/account/application/projection"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/contracts"
	sharedevents "wechat-clone/core/shared/contracts/events"
	infraMessaging "wechat-clone/core/shared/infra/messaging"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

//go:generate mockgen -package=processor -destination=processor_mock.go -source=processor.go
type Processor interface {
	Start() error
	Stop() error
}

type processor struct {
	consumer         []infraMessaging.Consumer
	accountReadRepo  accountprojection.AccountReadRepository
	searchProjection accountprojection.SearchProjection
}

func NewProcessor(cfg *config.Config, accountReadRepo accountprojection.AccountReadRepository, searchProjection accountprojection.SearchProjection) (Processor, error) {
	instance := &processor{
		consumer:         make([]infraMessaging.Consumer, 0, 1),
		accountReadRepo:  accountReadRepo,
		searchProjection: searchProjection,
	}

	topic := strings.TrimSpace(cfg.KafkaConfig.KafkaAccountConsumer.AccountOutboxTopic)
	if topic == "" || accountReadRepo == nil || searchProjection == nil {
		return instance, nil
	}

	handlerName := fmt.Sprintf("account-projection-%s-handler", strings.ToLower(topic))
	consumer, err := infraMessaging.NewConsumer(&infraMessaging.Config{
		Servers:      cfg.KafkaConfig.KafkaServers,
		Group:        cfg.KafkaConfig.KafkaAccountConsumer.AccountProjectionGroup,
		OffsetReset:  cfg.KafkaConfig.KafkaOffsetReset,
		ConsumeTopic: []string{topic},
		HandlerName:  handlerName,
		DLQ:          true,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	consumer.SetHandler(func(ctx context.Context, value []byte) error {
		return instance.handleAccountOutboxEvent(ctx, value)
	})
	instance.consumer = append(instance.consumer, consumer)

	return instance, nil
}

func (p *processor) Start() error {
	for _, consumer := range p.consumer {
		consumer.Read(infraMessaging.WrapConsumerCallback(consumer, "Handle account projection message failed"))
	}
	return nil
}

func (p *processor) Stop() error {
	infraMessaging.StopConsumers(p.consumer)
	return nil
}

func (p *processor) handleAccountOutboxEvent(ctx context.Context, value []byte) error {
	log := logging.FromContext(ctx).Named("AccountProjection")

	var event contracts.OutboxMessage
	if err := json.Unmarshal(value, &event); err != nil {
		return stackErr.Error(fmt.Errorf("unmarshal account outbox event failed: %w", err))
	}

	log.Infow("handle account outbox event",
		zap.String("event_name", event.EventName),
		zap.String("aggregate_id", event.AggregateID),
	)

	switch event.EventName {
	case sharedevents.EventAccountCreated,
		sharedevents.EventAccountUpdated,
		sharedevents.EventAccountProfileUpdated,
		sharedevents.EventAccountEmailVerified,
		sharedevents.EventAccountPasswordChanged,
		sharedevents.EventAccountBanned:
		return p.syncAccount(ctx, event.AggregateID)
	default:
		return nil
	}
}

func (p *processor) syncAccount(ctx context.Context, accountID string) error {
	if strings.TrimSpace(accountID) == "" {
		return stackErr.Error(fmt.Errorf("account projection aggregate_id is empty"))
	}
	if p.accountReadRepo == nil || p.searchProjection == nil {
		return nil
	}

	account, err := p.accountReadRepo.GetAccountByID(ctx, accountID)
	if err != nil {
		return stackErr.Error(fmt.Errorf("load account read model for projection failed: %w", err))
	}
	if err := p.searchProjection.SyncAccount(ctx, account); err != nil {
		return stackErr.Error(fmt.Errorf("sync account search projection failed: %w", err))
	}
	return nil
}
