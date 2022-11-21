package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
)

func ValidatorFunc(c *gin.Context, message string, code int) {
  valErr := api.Error{
    Code: 400,
    Message: message,
  }

  c.AbortWithStatusJSON(400, valErr)
}

