package middleware

import (
	"context"
	appCtx "go-socket/core/context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func AuthenMiddleware(appCtx *appCtx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" && websocket.IsWebSocketUpgrade(c.Request) {
			token = c.Query("token")
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		claims, err := appCtx.GetPaseto().ParseToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx := context.WithValue(c.Request.Context(), "account", claims)
		c.Request = c.Request.WithContext(ctx)
		c.Set("account", claims)
		c.Next()
	}
}
