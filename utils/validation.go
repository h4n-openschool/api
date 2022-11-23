package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
)

func ValidatorFunc(c *gin.Context, message string, code int) {
	if message == "no matching operation was found" {
		code = 404
	}

	valErr := api.Error{
		Code:    code,
		Message: message,
	}

	c.AbortWithStatusJSON(400, valErr)
}
