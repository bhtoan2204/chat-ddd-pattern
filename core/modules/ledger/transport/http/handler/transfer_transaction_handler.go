// CODE_GENERATOR - do not edit: handler
package handler

import (
	"net/http"

	"go-socket/core/modules/ledger/application/dto/in"
	"go-socket/core/modules/ledger/application/dto/out"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type transferTransactionHandler struct {
	transferTransaction cqrs.Dispatcher[*in.TransferTransactionRequest, *out.TransactionTransactionResponse]
}

func NewTransferTransactionHandler(
	transferTransaction cqrs.Dispatcher[*in.TransferTransactionRequest, *out.TransactionTransactionResponse],
) *transferTransactionHandler {
	return &transferTransactionHandler{
		transferTransaction: transferTransaction,
	}
}

func (h *transferTransactionHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.TransferTransactionRequest
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

	result, err := h.transferTransaction.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("TransferTransaction failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
