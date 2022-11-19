package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrValidation(ctx *gin.Context, err error) {
	var valErrs validator.ValidationErrors
	var ok bool
	if valErrs, ok = err.(validator.ValidationErrors); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	errs := []gin.H{}

	for _, valErr := range valErrs {
		errs = append(errs, gin.H{
			"field": valErr.Field(),
			"error": valErr.Error(),
		})
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "fields": errs})
}
