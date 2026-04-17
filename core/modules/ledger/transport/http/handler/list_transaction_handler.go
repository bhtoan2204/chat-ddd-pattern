// CODE_GENERATOR - do not edit: handler
package handler

import (
	"net/http"

	"wechat-clone/core/modules/ledger/application/dto/in"
	"wechat-clone/core/modules/ledger/application/dto/out"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type listTransactionHandler struct {
	listTransaction cqrs.Dispatcher[*in.ListTransactionRequest, *out.ListTransactionResponse]
}

func NewListTransactionHandler(
	listTransaction cqrs.Dispatcher[*in.ListTransactionRequest, *out.ListTransactionResponse],
) *listTransactionHandler {
	return &listTransactionHandler{
		listTransaction: listTransaction,
	}
}

func (h *listTransactionHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.ListTransactionRequest
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

	result, err := h.listTransaction.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("ListTransaction failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
