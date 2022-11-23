package utils

import (
	"time"

	middleware "github.com/deepmap/oapi-codegen/pkg/gin-middleware"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
	"go.uber.org/zap"
)

func ApplyMiddleware(e *gin.Engine, logger *zap.Logger) *gin.Engine {
	e = applyCorsMiddleware(e)
	e = applyValidationMiddleware(e)

	// Add error handling middleware to catch, log, and respond to errors.
	e.Use(ErrorHandler(logger))

	// Configure logging and recovery through Zap logger
	e.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	e.Use(ginzap.RecoveryWithZap(logger, true))

	return e
}

// applyCorsMiddleware applies the CORS middleware from gin-contrib, allowing
// all origins to query the server.
func applyCorsMiddleware(e *gin.Engine) *gin.Engine {
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true

	e.Use(cors.New(corsConf))
	return e
}

// applyValidationMiddleware adds request validation as generated by the OpenAPI
// code generation.
func applyValidationMiddleware(e *gin.Engine) *gin.Engine {
	// Add request validation from codegen
	swagger, _ := api.GetSwagger()
	opts := middleware.Options{
		ErrorHandler: ValidatorFunc,
	}

	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, &opts))
	return e
}
