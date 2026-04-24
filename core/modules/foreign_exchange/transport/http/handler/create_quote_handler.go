// CODE_GENERATOR - do not edit: handler
package handler

import (
	"net/http"

	"wechat-clone/core/modules/foreign_exchange/application/dto/in"
	"wechat-clone/core/modules/foreign_exchange/application/dto/out"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createQuoteHandler struct {
	createQuote cqrs.Dispatcher[*in.CreateQuoteRequest, *out.CreateQuoteResponse]
}

func NewCreateQuoteHandler(
	createQuote cqrs.Dispatcher[*in.CreateQuoteRequest, *out.CreateQuoteResponse],
) *createQuoteHandler {
	return &createQuoteHandler{
		createQuote: createQuote,
	}
}

func (h *createQuoteHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.CreateQuoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.createQuote.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("CreateQuote failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
