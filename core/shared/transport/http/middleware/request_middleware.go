package middleware

import (
	"go-socket/core/shared/pkg/contxt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			uid, err := uuid.NewRandom()
			if err == nil {
				requestID = uid.String()
			}
		}

		ctx := contxt.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}
