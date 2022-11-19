package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/classes/api"
	"github.com/h4n-openschool/classes/bus"
	"github.com/h4n-openschool/classes/models"
	"github.com/h4n-openschool/classes/repos/classes"
	"github.com/h4n-openschool/classes/utils"
)

type OpenSchoolImpl struct {
	Repository classes.ClassRepository
	Bus        *bus.Bus
}

func (i *OpenSchoolImpl) ClassesList(ctx *gin.Context, params api.ClassesListParams) {
	pagination := utils.NewPaginationQuery()
  pagination.ReadFromOptional(params.Page, params.PerPage)

	classes, err := i.Repository.GetAll(pagination)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong. Ask your administrator to check the server logs.",
		})
		log.Println(err)
		return
	}
  total, err := i.Repository.Count()
  if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong. Ask your administrator to check the server logs.",
		})
		log.Println(err)
		return
  }

  paginationResponse := api.PaginationData{
    Total: total,
    Page: pagination.Page,
    PerPage: pagination.PerPage,
    FirstUrl: "/v1/classes?page=1",
    LastUrl: fmt.Sprintf("/v1/classes?page=%v", total / pagination.PerPage),
    NextUrl: fmt.Sprintf("/v1/classes?page=%v", pagination.Page + 1),
    PrevUrl: fmt.Sprintf("/v1/classes?page=%v", pagination.Page - 1),
  }

  classList := models.ClassesAsApiClassList(classes)

  response := api.ClassesListResponse{
    Classes: classList,
    Pagination: paginationResponse,
  }

	ctx.JSON(http.StatusOK, response)
}
