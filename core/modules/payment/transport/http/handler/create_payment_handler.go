// CODE_GENERATOR: handler
package handler

import (
	"net/http"

	"go-socket/core/modules/payment/application/dto/in"
	"go-socket/core/modules/payment/application/dto/out"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createPaymentHandler struct {
	createPayment cqrs.Dispatcher[*in.CreatePaymentRequest, *out.CreatePaymentResponse]
}

func NewCreatePaymentHandler(
	createPayment cqrs.Dispatcher[*in.CreatePaymentRequest, *out.CreatePaymentResponse],
) *createPaymentHandler {
	return &createPaymentHandler{
		createPayment: createPayment,
	}
}

func (h *createPaymentHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.CreatePaymentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, nil
	}

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, nil
	}

	if accountID, err := accountIDFromContext(ctx); err == nil {
		if request.DebitAccountID == "" {
			request.DebitAccountID = accountID
		} else if request.DebitAccountID != accountID {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "debit_account_id must match authenticated account"})
			return nil, nil
		}
	}

	result, err := h.createPayment.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("CreatePayment failed", zap.Error(err))
		writeProviderError(c, err)
		return nil, nil
	}

	c.JSON(http.StatusCreated, result)
	return nil, nil
}
