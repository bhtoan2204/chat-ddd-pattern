package middleware

import "github.com/gin-gonic/gin"

func AuthorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
