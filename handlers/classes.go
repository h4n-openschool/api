package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/h4n-openschool/api/api"
	"github.com/h4n-openschool/api/bus"
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/utils"
	"github.com/rabbitmq/amqp091-go"
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
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	// Get the total number of classes in the database to build the PaginationData object.
	total, err := i.Repository.Count()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
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

// ClassesCreate implements the classesCreate contract from the OpenAPI spec.
func (i *OpenSchoolImpl) ClassesCreate(ctx *gin.Context) {
	var body api.ClassesCreateRequest
	if err := ctx.BindJSON(&body); err != nil {
		_ = ctx.AbortWithError(400, err)
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
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	class, err := i.Repository.Create(in)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
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
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	response := api.ClassesCreateResponse{
		Class: class.AsApiClass(),
	}

	ctx.JSON(http.StatusCreated, response)
}

// ClassesGet implements the classesGet contract from the OpenAPI spec.
func (i *OpenSchoolImpl) ClassesGet(ctx *gin.Context, id api.Cuid) {
	class, err := i.Repository.Get(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if class == nil {
		_ = ctx.AbortWithError(http.StatusNotFound, errors.New("Class not found"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"class": class.AsApiClass()})
}

func (i *OpenSchoolImpl) ClassesUpdate(ctx *gin.Context, id api.Cuid) {
  var body api.ClassesUpdateJSONRequestBody
  _ = ctx.Bind(&body)

  class, err := i.Repository.Get(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if class == nil {
		_ = ctx.AbortWithError(http.StatusNotFound, errors.New("Class not found"))
		return
	}

  class = class.ReconcileWithApiClass(body.Name, body.Description, body.DisplayName)

  class, err = i.Repository.Update(class)
  if err != nil {
    _ = ctx.AbortWithError(http.StatusInternalServerError, err)
  }

  ctx.JSON(http.StatusOK, api.ClassesUpdateResponse{Class: class.AsApiClass()})
}
