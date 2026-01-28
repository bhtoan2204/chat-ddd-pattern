package middleware

import (
	"net/http"
	"strings"

	"go-socket/core/shared/infra/idempotency"

	"github.com/gin-gonic/gin"
)

const idempotencyHeader = "Idempotency-Key"

func IdempotencyMiddleware(manager *idempotency.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if manager == nil {
			c.Next()
			return
		}
		if !isWriteMethod(c.Request.Method) {
			c.Next()
			return
		}
		key := strings.TrimSpace(c.GetHeader(idempotencyHeader))
		if key == "" {
			c.Next()
			return
		}
		ok, err := manager.Begin(c.Request.Context(), key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "duplicate request"})
			return
		}
		c.Next()
		status := c.Writer.Status()
		success := status >= http.StatusOK && status < http.StatusMultipleChoices
		_ = manager.End(c.Request.Context(), key, success)
	}
}

func isWriteMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}
