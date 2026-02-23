package middleware

import (
	"net/http"
	"time"

	"go-socket/core/shared/infra/cache"
	"go-socket/core/shared/infra/ratelimit"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(cache cache.Cache) gin.HandlerFunc {
	limiter := ratelimit.NewSlidingWindowLimiter(cache, 60, time.Minute)
	return func(c *gin.Context) {
		if cache == nil {
			c.Next()
			return
		}
		ok, err := limiter.Allow(c.Request.Context(), c.ClientIP())
		if err != nil {
			c.Next()
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
