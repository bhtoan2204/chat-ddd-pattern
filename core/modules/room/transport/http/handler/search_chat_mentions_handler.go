package handler

import (
	"net/http"

	"wechat-clone/core/modules/room/application/dto/in"
	"wechat-clone/core/modules/room/application/dto/out"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type searchChatMentionsHandler struct {
	searchChatMentions cqrs.Dispatcher[*in.SearchChatMentionsRequest, []*out.ChatMentionCandidateResponse]
}

func NewSearchChatMentionsHandler(
	searchChatMentions cqrs.Dispatcher[*in.SearchChatMentionsRequest, []*out.ChatMentionCandidateResponse],
) *searchChatMentionsHandler {
	return &searchChatMentionsHandler{searchChatMentions: searchChatMentions}
}

func (h *searchChatMentionsHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)

	var request in.SearchChatMentionsRequest
	request.RoomID = c.Param("room_id")
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.searchChatMentions.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("SearchChatMentions failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
