package projection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-socket/core/modules/payment/domain/entity"
	sharedevents "go-socket/core/shared/contracts/events"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (p *processor) handleAccountEvent(ctx context.Context, value []byte) error {
	log := logging.FromContext(ctx).Named("PaymentAccountProjection")
	var event accountOutboxMessage
	if err := json.Unmarshal(value, &event); err != nil {
		return stackErr.Error(fmt.Errorf("unmarshal account outbox event failed: %v", err))
	}

	log.Infow("handle account event", zap.String("event_name", event.EventName))
	switch event.EventName {
	case "EventAccountCreated":
		return p.handleAccountCreatedEvent(ctx, &event)
	case "EventAccountUpdated":
		return p.handleAccountUpdatedEvent(ctx, &event)
	case "EventAccountBanned":
		return p.handleAccountBannedEvent(ctx, &event)
	default:
		return nil
	}
}

func (p *processor) handleAccountCreatedEvent(ctx context.Context, event *accountOutboxMessage) error {
	log := logging.FromContext(ctx).Named("handleAccountCreatedEvent")
	payload, err := decodeExternalEventPayload[sharedevents.AccountCreatedEvent](event.EventData)
	if err != nil {
		return stackErr.Error(fmt.Errorf("decode event payload failed: %v", err))
	}
	if payload == nil {
		return stackErr.Error(fmt.Errorf("invalid payload type for event %s", "EventAccountCreated"))
	}

	projection := &entity.PaymentAccount{
		ID:        payload.AccountID,
		AccountID: payload.AccountID,
		Email:     payload.Email,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.CreatedAt,
	}

	if err := p.accountProjectionRepo.UpsertAccountProjection(ctx, projection); err != nil {
		log.Errorw("upsert account projection failed", zap.Error(err))
		return stackErr.Error(fmt.Errorf("upsert account projection failed: %v", err))
	}
	return nil
}

func (p *processor) handleAccountUpdatedEvent(ctx context.Context, event *accountOutboxMessage) error {
	log := logging.FromContext(ctx).Named("handleAccountUpdatedEvent")
	payload, err := decodeExternalEventPayload[sharedevents.AccountUpdatedEvent](event.EventData)
	if err != nil {
		return stackErr.Error(fmt.Errorf("decode event payload failed: %v", err))
	}
	if payload == nil {
		return stackErr.Error(fmt.Errorf("invalid payload type for event %s", "EventAccountUpdated"))
	}

	existing, err := p.accountProjectionRepo.GetAccountProjectionByAccountID(ctx, payload.AccountID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorw("get account projection failed", zap.Error(err))
		return stackErr.Error(fmt.Errorf("get account projection failed: %v", err))
	}

	if existing == nil {
		projection := &entity.PaymentAccount{
			ID:        payload.AccountID,
			AccountID: payload.AccountID,
			Email:     payload.Email,
			CreatedAt: payload.UpdatedAt,
			UpdatedAt: payload.UpdatedAt,
		}
		if err := p.accountProjectionRepo.CreateAccountProjection(ctx, projection); err != nil {
			log.Errorw("create account projection failed", zap.Error(err))
			return stackErr.Error(fmt.Errorf("create account projection failed: %v", err))
		}
		return nil
	}

	existing.Email = payload.Email
	existing.AccountID = payload.AccountID
	existing.UpdatedAt = payload.UpdatedAt
	if err := p.accountProjectionRepo.UpdateAccountProjection(ctx, existing); err != nil {
		log.Errorw("update account projection failed", zap.Error(err))
		return stackErr.Error(fmt.Errorf("update account projection failed: %v", err))
	}
	return nil
}

func (p *processor) handleAccountBannedEvent(ctx context.Context, event *accountOutboxMessage) error {
	log := logging.FromContext(ctx).Named("handleAccountBannedEvent")
	payload, err := decodeExternalEventPayload[sharedevents.AccountBannedEvent](event.EventData)
	if err != nil {
		return stackErr.Error(fmt.Errorf("decode event payload failed: %v", err))
	}
	if payload == nil {
		return stackErr.Error(fmt.Errorf("invalid payload type for event %s", "EventAccountBanned"))
	}

	if err := p.accountProjectionRepo.DeleteAccountProjection(ctx, payload.AccountID); err != nil {
		log.Errorw("delete account projection failed", zap.Error(err))
		return stackErr.Error(fmt.Errorf("delete account projection failed: %v", err))
	}
	return nil
}
