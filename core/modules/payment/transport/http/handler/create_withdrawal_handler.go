// CODE_GENERATOR - do not edit: handler
package handler

import (
	"net/http"

	"wechat-clone/core/modules/payment/application/dto/in"
	"wechat-clone/core/modules/payment/application/dto/out"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createWithdrawalHandler struct {
	createWithdrawal cqrs.Dispatcher[*in.CreateWithdrawalRequest, *out.CreateWithdrawalResponse]
}

func NewCreateWithdrawalHandler(
	createWithdrawal cqrs.Dispatcher[*in.CreateWithdrawalRequest, *out.CreateWithdrawalResponse],
) *createWithdrawalHandler {
	return &createWithdrawalHandler{
		createWithdrawal: createWithdrawal,
	}
}

func (h *createWithdrawalHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.CreateWithdrawalRequest
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

	result, err := h.createWithdrawal.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("CreateWithdrawal failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	c.JSON(201, result)
	return nil, nil
}
