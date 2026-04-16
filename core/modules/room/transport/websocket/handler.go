package socket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	appCtx "go-socket/core/context"
	"go-socket/core/modules/room/constant"
	"go-socket/core/modules/room/types"
	"go-socket/core/shared/pkg/actorctx"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/pubsub"
	"go-socket/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type wsHandler struct {
	hub       IHub
	upgrader  websocket.Upgrader
	subcriber *pubsub.Subscription
}

func NewWSHandler(appContext *appCtx.AppContext, hub IHub, upgrader websocket.Upgrader) *wsHandler {
	log := logging.FromContext(context.Background())
	subcriber, err := appContext.LocalBus().Subscribe(constant.RealtimeMessageTopic)
	if err != nil {
		return nil
	}

	go func() {
		for msg := range subcriber.C() {
			if err := handleRealtimeMessage(context.Background(), hub, msg); err != nil {
				log.Warnw("handle realtime message failed", zap.Error(err), zap.Any("msg", msg))
			}
		}
		log.Warnw("channel closed unexpectedly")
	}()

	return &wsHandler{
		hub:       hub,
		upgrader:  upgrader,
		subcriber: subcriber,
	}
}

func (h *wsHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()
	log := logging.FromContext(ctx)

	if h.hub == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "websocket hub is not initialized"})
		return
	}

	accountID, err := actorctx.AccountIDFromContext(ctx)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorw("failed to upgrade websocket connection", zap.Error(err))
		return
	}

	client := NewClient(ctx, conn, c.Query("client_id"), accountID)
	h.hub.Register(ctx, client)

	clientCtx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		client.ReadPump(clientCtx, h.hub)
	}()

	go func() {
		defer cancel()
		client.WritePump(clientCtx)
	}()
}

func handleRealtimeMessage(ctx context.Context, hub IHub, msg pubsub.Message) error {
	if hub == nil {
		return stackErr.Error(errors.New("hub is nil"))
	}

	switch msg.Topic {
	case constant.RealtimeMessageTopic:
		event, err := asRealtimeChatMessageCreatedEvent(msg.Data)
		if err != nil {
			return stackErr.Error(err)
		}
		payload, err := json.Marshal(event.Payload)
		if err != nil {
			return stackErr.Error(fmt.Errorf("marshal realtime chat message payload: %v", err))
		}

		if err := hub.Publish(ctx, Message{
			RoomID:       event.RoomId,
			Action:       event.Type,
			Data:         payload,
			RecipientIDs: event.RecipientIds,
		}); err != nil {
			return stackErr.Error(err)
		}
		return nil

	default:
		return nil
	}
}

func asRealtimeChatMessageCreatedEvent(data any) (types.MessagePayload, error) {
	v, ok := data.(types.MessagePayload)
	if !ok {
		return types.MessagePayload{}, fmt.Errorf("invalid payload type: %T", data)
	}
	return v, nil
}
