package handler

import (
	"context"
	"errors"
	"net/http"

	paymentservice "go-socket/core/modules/payment/application/service"
	"go-socket/core/modules/payment/providers"
	"go-socket/core/shared/infra/xpaseto"

	"github.com/gin-gonic/gin"
)

func accountIDFromContext(ctx context.Context) (string, error) {
	account, ok := ctx.Value("account").(*xpaseto.PasetoPayload)
	if !ok || account == nil || account.AccountID == "" {
		return "", errors.New("account not found")
	}
	return account.AccountID, nil
}

func writeProviderError(c *gin.Context, err error) {
	status := http.StatusInternalServerError

	switch {
	case isProviderValidation(err):
		status = http.StatusBadRequest
	case isProviderDuplicate(err):
		status = http.StatusConflict
	case isProviderNotFound(err):
		status = http.StatusNotFound
	case isUnknownProvider(err):
		status = http.StatusBadRequest
	case errors.Is(err, providers.ErrInvalidWebhookSignature):
		status = http.StatusUnauthorized
	}

	c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
}

func isProviderValidation(err error) bool {
	return errors.Is(err, paymentservice.ErrValidation)
}

func isProviderDuplicate(err error) bool {
	return errors.Is(err, paymentservice.ErrDuplicateTransaction) || errors.Is(err, paymentservice.ErrDuplicatePayment)
}

func isProviderNotFound(err error) bool {
	return errors.Is(err, paymentservice.ErrTransactionNotFound) || errors.Is(err, paymentservice.ErrPaymentIntentNotFound)
}

func isUnknownProvider(err error) bool {
	return errors.Is(err, providers.ErrProviderNotFound)
}
