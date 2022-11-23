package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
	"github.com/h4n-openschool/api/bus"
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/repos/teachers"
	"github.com/h4n-openschool/api/utils"
	"github.com/rabbitmq/amqp091-go"
)

// Teachers implements the teachersList operation from the OpenAPI
// specification in [../api/spec.yaml].
func (i *OpenSchoolImpl) TeachersList(ctx *gin.Context, params api.TeachersListParams) {
	// Read pagination options from the ClassesListParams object
	pagination := utils.NewPaginationQuery()
	pagination.ReadFromOptional(params.Page, params.PerPage)

	// Retrieve a paginated list of classes
	classes, err := i.TeacherRepository.GetAll(pagination)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	// Get the total number of classes in the database to build the PaginationData object.
	total, err := i.ClassRepository.Count()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	// Generate pagination data from the total and input pagination options.
	paginationData := utils.GeneratePaginationData("/v1/teachers", total, pagination)

	// Convert the class model array to an api.ClassList type to meet the OpenAPI definition.
	teachersList := models.TeachersAsApiTeacherList(classes)

	// Build the response body as a ClassesListResponse type.
	response := api.TeachersListResponse{
		Teachers:   teachersList,
		Pagination: paginationData,
	}

	// Respond to the request in JSON
	ctx.JSON(http.StatusOK, response)
}

// TeachersCreate implements the teachersCreate contract from the OpenAPI spec.
func (i *OpenSchoolImpl) TeachersCreate(ctx *gin.Context) {
	var body api.TeachersCreateJSONRequestBody
	_ = ctx.Bind(&body)

	in := models.Teacher{
		FullName: body.FullName,
		Email:    body.Email,
	}

	c, q, err := i.Bus.ChannelQueue()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	teacher, err := i.TeacherRepository.Create(in)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ev := bus.Event[*models.Teacher]{
		Data: teacher,
		Metadata: bus.EventMeta{
			Type:     "create",
			Resource: "teacher",
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

	response := api.TeachersCreateResponse{
		Teacher: teacher.AsApiTeacher(),
	}

	ctx.JSON(http.StatusCreated, response)
}

// TeachersGet implements the teachersGet contract from the OpenAPI spec.
func (i *OpenSchoolImpl) TeachersGet(ctx *gin.Context, id api.Cuid) {
	teacher, err := i.TeacherRepository.Get(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if teacher == nil {
		_ = ctx.AbortWithError(http.StatusNotFound, errors.New("Teacher not found"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"teacher": teacher.AsApiTeacher()})
}

// TeachersUpdate implements the teachersUpdate contract from the OpenAPI spec.
func (i *OpenSchoolImpl) TeachersUpdate(ctx *gin.Context, id api.Cuid) {
	var body api.TeachersUpdateJSONRequestBody
	_ = ctx.Bind(&body)

	teacher := &models.Teacher{}
	teacher.Id = id
	teacher.FullName = body.FullName
	teacher.Email = body.Email

	teacher, err := i.TeacherRepository.Update(teacher)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, api.TeachersUpdateResponse{Teacher: teacher.AsApiTeacher()})
}

// TeachersDelete implements the teachersDelete contract from the OpenAPI spec.
func (i *OpenSchoolImpl) TeachersDelete(ctx *gin.Context, id api.Cuid) {
	teacher := models.Teacher{}
	teacher.Id = id

	err := i.TeacherRepository.Delete(teacher)
	if err != nil {
		if err == teachers.TeacherDoesNotExist {
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
