package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Wrap(h interface {
	Handle(c *gin.Context) (interface{}, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if h == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "handler is nil"})
			return
		}
		data, err := h.Handle(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	}
}
