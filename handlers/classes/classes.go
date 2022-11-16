package classes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	repos "github.com/h4n-openschool/classes/repos/classes"
	"github.com/h4n-openschool/classes/utils"
	"github.com/lucsky/cuid"
)

// @BasePath /classes

// GetClass godoc
// @Summary Get a class by ID
// @Schemes
// @Tags classes
// @Produce json
// @Success 200
// @Router /classes/:id [get]
func GetClass(repo repos.ClassRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
    // Get the ID from the route parameter
		id := ctx.Param("id")

    // Validate the ID is a CUID
		if err := cuid.IsCuid(id); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID, must be a cuid.",
			})
			return
		}

    // Get a class by its CUID (passed in route params)
		class, err := repo.Get(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong. Ask your administrator to check the server logs.",
			})
			log.Println(err)
			return
		}

    // If the class doesn't exist, 404
		if class == nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "No such class exists.",
			})
			return
		}

    // Respond with the class
		ctx.JSON(http.StatusOK, gin.H{"class": class})
	}
}

// GetClasses godoc
// @Summary Get a paginated list of classes
// @Schemes
// @Tags classes
// @Produce json
// @Success 200
// @Router /classes [get]
func GetClasses(repo repos.ClassRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pagination := utils.NewPaginationQuery()
		pagination.Read(ctx)

		classes, err := repo.GetAll(pagination)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong. Ask your administrator to check the server logs.",
			})
			log.Println(err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"classes": classes, "pagination": pagination})
	}
}
