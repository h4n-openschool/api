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
			status := -1
			message := err.Error()

			logger.Info(message)

			if err.Error() == "no matching operation was found" {
				status = 404
				message = "Route not found."
			}

			e := api.Error{
				Code:    status,
				Message: message,
			}

			c.JSON(status, e)
		}
	}
}
