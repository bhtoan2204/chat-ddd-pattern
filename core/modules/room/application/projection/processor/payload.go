package processor

import (
	"context"
	"fmt"
	"reflect"

	roomprojection "wechat-clone/core/modules/room/application/projection"
	"wechat-clone/core/shared/contracts"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

var projectionEventPayloadTypes = map[string]reflect.Type{
	roomprojection.EventRoomAggregateProjectionSynced:    reflect.TypeOf(roomprojection.RoomAggregateSync{}),
	roomprojection.EventRoomAggregateProjectionDeleted:   reflect.TypeOf(roomprojection.RoomAggregateDeleted{}),
	roomprojection.EventMessageAggregateProjectionSynced: reflect.TypeOf(roomprojection.MessageAggregateSync{}),
}

func decodeEventPayload(ctx context.Context, eventName string, raw []byte) (interface{}, error) {
	logger := logging.FromContext(ctx)
	payloadType, ok := projectionEventPayloadTypes[eventName]
	if !ok {
		logger.Warnw("unsupported event_name", zap.String("event_name", eventName))
		return nil, nil
	}

	payload := reflect.New(payloadType).Interface()
	if err := contracts.UnmarshalEventData(raw, payload); err != nil {
		logger.Errorw("unmarshal event_data failed", zap.Error(err), zap.String("raw", string(raw)))
		return nil, stackErr.Error(fmt.Errorf("unmarshal event_data failed: %w", err))
	}

	return payload, nil
}
