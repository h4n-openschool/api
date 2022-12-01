package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/h4n-openschool/api/api"
	"github.com/h4n-openschool/api/auth"
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/repos/classes"
	"github.com/h4n-openschool/api/utils"
)

// ClassesList implements the classesList operation from the OpenAPI
// specification in [../api/spec.yaml].
func (i *OpenSchoolImpl) ClassesList(ctx *gin.Context, params api.ClassesListParams) {
	if ok := auth.MustAuthenticate(ctx, i.TeacherRepository); ok {
		return
	}

	// Read pagination options from the ClassesListParams object
	pagination := utils.NewPaginationQuery()
	pagination.ReadFromOptional(params.Page, params.PerPage)

	// Retrieve a paginated list of classes
	classes, err := i.ClassRepository.GetAll(pagination)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	// Get the total number of classes in the database to build the PaginationData object.
	total, err := i.ClassRepository.Count()
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

	sd, err := time.Parse(time.RFC3339, *body.StartDate)
	if err != nil {
		_ = ctx.AbortWithError(400, err)
	}

	ed, err := time.Parse(time.RFC3339, *body.EndDate)
	if err != nil {
		_ = ctx.AbortWithError(400, err)
	}

	in := models.Class{
		DisplayName: body.DisplayName,
		StartDate:   sd,
		EndDate:     ed,
	}

	if body.Name != nil {
		in.Name = *body.Name
	} else {
		in.Name = slug.Make(body.DisplayName)
	}

	if body.Description != nil {
		in.Description = body.Description
	}

	class, err := i.ClassRepository.Create(in)
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
	class, err := i.ClassRepository.Get(id)
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

	sd, err := time.Parse(time.RFC3339, *body.StartDate)
	if err != nil {
		_ = ctx.AbortWithError(400, err)
	}

	ed, err := time.Parse(time.RFC3339, *body.EndDate)
	if err != nil {
		_ = ctx.AbortWithError(400, err)
	}

	class := &models.Class{
    StartDate: sd,
    EndDate: ed,
  }
	class.Id = id
	class = class.ReconcileWithApiClass(body.Description, body.DisplayName)

	class, err = i.ClassRepository.Update(class)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, api.ClassesUpdateResponse{Class: class.AsApiClass()})
}

func (i *OpenSchoolImpl) ClassesDelete(ctx *gin.Context, id api.Cuid) {
	class := models.Class{}
	class.Id = id

	err := i.ClassRepository.Delete(class)
	if err != nil {
		if err == classes.ClassDoesNotExist {
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
