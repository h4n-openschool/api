package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/classes/repos"
	"github.com/lucsky/cuid"
)

func GetClass(repo repos.ClassRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := cuid.IsCuid(id); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID, must be a cuid.",
			})
			return
		}

		class, err := repo.Get(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong. Ask your administrator to check the server logs.",
			})
			log.Println(err)
			return
		}

		if class == nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "No such class exists.",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"class": class})
	}
}
