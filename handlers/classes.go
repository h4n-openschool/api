package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/classes/api"
	"github.com/h4n-openschool/classes/bus"
	"github.com/h4n-openschool/classes/models"
	"github.com/h4n-openschool/classes/utils"
	"github.com/rabbitmq/amqp091-go"
  "github.com/gosimple/slug"
)

// ClassesList implements the classesList operation from the OpenAPI
// specification in [../api/spec.yaml].
func (i *OpenSchoolImpl) ClassesList(ctx *gin.Context, params api.ClassesListParams) {
	// Read pagination options from the ClassesListParams object
	pagination := utils.NewPaginationQuery()
	pagination.ReadFromOptional(params.Page, params.PerPage)

	// Retrieve a paginated list of classes
	classes, err := i.Repository.GetAll(pagination)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong. Ask your administrator to check the server logs.",
		})
		log.Println(err)
		return
	}

	// Get the total number of classes in the database to build the PaginationData object.
	total, err := i.Repository.Count()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong. Ask your administrator to check the server logs.",
		})
		log.Println(err)
		return
	}

	// Generate pagination data from the total and input pagination options.
	paginationData := utils.GeneratePaginationData("/v1/classes", total, pagination)

	// Convert the class model array to an api.ClassList type to meet the OpenAPI definition.
	classList := models.ClassesAsApiClassList(classes)

	// Build the response body as a ClassesListResponse type.
	response := api.ClassesListResponse{
		Classes:    classList,
		Pagination: paginationData,
	}

	// Respond to the request in JSON
	ctx.JSON(http.StatusOK, response)
}

func (i *OpenSchoolImpl) ClassesCreate(ctx *gin.Context) {
  var body api.ClassesCreateRequest
  if err := ctx.BindJSON(&body); err != nil {
    utils.ErrValidation(ctx, err)
    return
  }

  in := models.Class{
    DisplayName: body.DisplayName,
  }

  if body.Name != nil {
    in.Name = *body.Name
  } else {
    in.Name = slug.Make(body.DisplayName)
  }

  if body.Description != nil {
    in.Description = body.Description
  }

  c, q, err := i.Bus.ChannelQueue()
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }

  class, err := i.Repository.Create(in)
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }

  ev := bus.Event[*models.Class]{
    Data: class,
    Metadata: bus.EventMeta{
      Type:     "create",
      Resource: "class",
    },
  }

  j, _ := json.Marshal(ev)

  err = c.PublishWithContext(
    ctx.Request.Context(),
    "",
    q.Name,
    false,
    false,
    amqp091.Publishing{
      ContentType: "application/json",
      Body:        j,
    },
  )
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }

  response := api.ClassesCreateResponse{
    Class: class.AsApiClass(),
  }

  ctx.JSON(http.StatusCreated, response)
}
