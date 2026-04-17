package middleware

import (
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func NewAppError(code int, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     stackErr.Error(err),
	}
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		if appErr, ok := err.(*AppError); ok {
			c.JSON(appErr.Code, gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})

			return
		}

		c.JSON(500, gin.H{
			"code":    500,
			"message": "Internal server error",
		})
	}
}
