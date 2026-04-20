package repository

import (
	"context"
	"errors"
	"fmt"
	"wechat-clone/core/modules/relationship/domain/aggregate"
	"wechat-clone/core/modules/relationship/domain/entity"
	"wechat-clone/core/modules/relationship/domain/repos"
	"wechat-clone/core/modules/relationship/infra/persistent/models"

	eventpkg "wechat-clone/core/shared/pkg/event"
	"wechat-clone/core/shared/pkg/stackErr"

	"gorm.io/gorm"
)

type friendRequestAggregateRepo struct {
	db         *gorm.DB
	serializer eventpkg.Serializer
}

func newFriendRequestAggregateRepo(db *gorm.DB) repos.FriendRequestAggregateRepository {
	serializer := eventpkg.NewSerializer()

	return &friendRequestAggregateRepo{
		db:         db,
		serializer: serializer,
	}
}

func (f *friendRequestAggregateRepo) Load(ctx context.Context, friendRequestID string) (*aggregate.FriendRequestAggregate, error) {
	if friendRequestID == "" {
		return nil, stackErr.Error(fmt.Errorf("friend request id is required"))
	}

	var friendRequestModel models.FriendRequest
	err := f.db.WithContext(ctx).
		Where("id = ?", friendRequestID).
		First(&friendRequestModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, stackErr.Error(err)
		}
		return nil, stackErr.Error(err)
	}

	friendRequestEntity := toFriendRequestEntity(&friendRequestModel)

	agg, err := aggregate.NewFriendRequest(friendRequestID)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	if err := agg.SetFriendRequest(friendRequestEntity); err != nil {
		return nil, stackErr.Error(err)
	}

	return agg, nil
}

func (f *friendRequestAggregateRepo) Save(ctx context.Context, agg *aggregate.FriendRequestAggregate) error {
	if agg == nil {
		return stackErr.Error(fmt.Errorf("friend request aggregate is required"))
	}
	if agg.FriendRequest == nil {
		return stackErr.Error(fmt.Errorf("friend request entity is required"))
	}

	aggID := agg.AggregateID()
	if aggID == "" {
		return stackErr.Error(fmt.Errorf("friend request aggregate id is required"))
	}

	agg.FriendRequest.ID = aggID

	friendRequestModel := toFriendRequestModel(agg.FriendRequest)
	if friendRequestModel == nil {
		return stackErr.Error(fmt.Errorf("friend request model is nil"))
	}

	return f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(friendRequestModel).Error; err != nil {
			return stackErr.Error(err)
		}

		events := agg.Events()
		if len(events) > 0 {
			outboxModels := make([]models.RelationOutboxEvent, 0, len(events))

			for _, ev := range events {
				eventData, err := f.serializeEventData(ev.EventData)
				if err != nil {
					return stackErr.Error(err)
				}

				outboxModels = append(outboxModels, models.RelationOutboxEvent{
					AggregateID:   ev.AggregateID,
					AggregateType: ev.AggregateType,
					Version:       int64(ev.Version),
					EventName:     ev.EventName,
					EventData:     eventData,
				})
			}

			if err := tx.Create(&outboxModels).Error; err != nil {
				return stackErr.Error(err)
			}
		}

		agg.Update()
		return nil
	})
}

func (f *friendRequestAggregateRepo) serializeEventData(data interface{}) (string, error) {
	if data == nil {
		return "", stackErr.Error(fmt.Errorf("event data is nil"))
	}

	switch v := data.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		b, err := f.serializer.Marshal(data)
		if err != nil {
			return "", stackErr.Error(err)
		}
		return string(b), nil
	}
}

func toFriendRequestModel(e *entity.FriendRequest) *models.FriendRequest {
	if e == nil {
		return nil
	}

	return &models.FriendRequest{
		ID:             e.ID,
		RequesterID:    e.RequesterID,
		AddresseeID:    e.AddresseeID,
		Status:         toFriendRequestModelStatus(e.Status),
		Message:        e.Message,
		CreatedAt:      e.CreatedAt,
		RespondedAt:    e.RespondedAt,
		ExpiredAt:      e.ExpiredAt,
		CancelledAt:    e.CancelledAt,
		RejectedReason: e.RejectedReason,
	}
}

func toFriendRequestEntity(m *models.FriendRequest) *entity.FriendRequest {
	if m == nil {
		return nil
	}

	return &entity.FriendRequest{
		ID:             m.ID,
		RequesterID:    m.RequesterID,
		AddresseeID:    m.AddresseeID,
		Status:         toFriendRequestEntityStatus(m.Status),
		Message:        m.Message,
		CreatedAt:      m.CreatedAt,
		RespondedAt:    m.RespondedAt,
		ExpiredAt:      m.ExpiredAt,
		CancelledAt:    m.CancelledAt,
		RejectedReason: m.RejectedReason,
	}
}

func toFriendRequestModels(es []*entity.FriendRequest) []*models.FriendRequest {
	if es == nil {
		return nil
	}

	out := make([]*models.FriendRequest, 0, len(es))
	for _, e := range es {
		out = append(out, toFriendRequestModel(e))
	}
	return out
}

func toFriendRequestEntities(ms []*models.FriendRequest) []*entity.FriendRequest {
	if ms == nil {
		return nil
	}

	out := make([]*entity.FriendRequest, 0, len(ms))
	for _, m := range ms {
		out = append(out, toFriendRequestEntity(m))
	}
	return out
}

func toFriendRequestModelStatus(s entity.FriendRequestStatus) models.FriendRequestStatus {
	return models.FriendRequestStatus(s)
}

func toFriendRequestEntityStatus(s models.FriendRequestStatus) entity.FriendRequestStatus {
	return entity.FriendRequestStatus(s)
}
