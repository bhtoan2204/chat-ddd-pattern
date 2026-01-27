package middleware

import (
	appCtx "go-socket/core/context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenMiddleware(appCtx *appCtx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
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
		c.Set("account", claims)
		c.Next()
	}
}
