package middleware

import (
	"context"
	"net/http"
	"strings"
	appCtx "wechat-clone/core/context"
	"wechat-clone/core/shared/pkg/actorctx"

	"github.com/gin-gonic/gin"
)

func extractToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token != "" {
		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
		if token != "" {
			return token
		}
	}

	token = strings.TrimSpace(c.Query("authorization"))
	if token != "" {
		return token
	}

	return ""
}

type contextKey string

const accountContextKey contextKey = "account"

func AuthenMiddleware(appCtx *appCtx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		claims, err := appCtx.GetPaseto().ParseAccessToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx := actorctx.WithActor(c.Request.Context(), actorctx.Actor{
			AccountID: claims.AccountID,
			Email:     claims.Email,
		})
		ctx = context.WithValue(ctx, accountContextKey, claims)
		c.Request = c.Request.WithContext(ctx)
		c.Set(accountContextKey, claims)
		c.Next()
	}
}
