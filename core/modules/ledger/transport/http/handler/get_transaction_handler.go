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

type getTransactionHandler struct {
	getTransaction cqrs.Dispatcher[*in.GetTransactionRequest, *out.TransactionResponse]
}

func NewGetTransactionHandler(
	getTransaction cqrs.Dispatcher[*in.GetTransactionRequest, *out.TransactionResponse],
) *getTransactionHandler {
	return &getTransactionHandler{
		getTransaction: getTransaction,
	}
}

func (h *getTransactionHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.GetTransactionRequest
	request.TransactionID = c.Param("transaction_id")

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.getTransaction.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("GetTransaction failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
