package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
	"go.uber.org/zap"
)

// ErrorHandler responds errors in a user-friendly format and logs them to the
// console.
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			e := api.Error{
				Code:    c.Writer.Status(),
				Message: err.Error(),
			}

			c.JSON(-1, e)
		}
	}
}
