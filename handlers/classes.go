package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/classes/models"
	"github.com/h4n-openschool/classes/repos"
	"github.com/h4n-openschool/classes/utils"
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

type CreateClassDto struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description,omitempty"`
}

func CreateClass(repo repos.ClassRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := CreateClassDto{}
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		in := models.Class{
			Name:        body.Name,
			DisplayName: body.DisplayName,
			Description: body.Description,
		}

		class, err := repo.Create(in)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"class": class,
		})
	}
}
