package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
	"github.com/h4n-openschool/api/auth"
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/repos/students"
	"github.com/h4n-openschool/api/utils"
)

// StudentsList implements the studentsList operation from the OpenAPI
// specification in [../api/spec.yaml].
func (i *OpenSchoolImpl) StudentsList(ctx *gin.Context, params api.StudentsListParams) {
  if ok := auth.MustAuthenticate(ctx, i.TeacherRepository); ok {
    return
  }

	// Read pagination options from the StudentsListParams object
	pagination := utils.NewPaginationQuery()
	pagination.ReadFromOptional(params.Page, params.PerPage)

	// Retrieve a paginated list of students
	students, err := i.StudentRepository.GetAll(pagination)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	// Get the total number of students in the database to build the PaginationData object.
	total, err := i.StudentRepository.Count()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	// Generate pagination data from the total and input pagination options.
	paginationData := utils.GeneratePaginationData("/v1/students", total, pagination)

	// Convert the student model array to an api.StudentList type to meet the OpenAPI definition.
	studentList := models.StudentsAsApiStudentList(students)

	// Build the response body as a StudentsListResponse type.
	response := api.StudentsListResponse{
		Students:    studentList,
		Pagination: paginationData,
	}

	// Respond to the request in JSON
	ctx.JSON(http.StatusOK, response)
}

// StudentsCreate implements the studentsCreate contract from the OpenAPI spec.
func (i *OpenSchoolImpl) StudentsCreate(ctx *gin.Context) {
	var body api.StudentsCreateRequest
	if err := ctx.BindJSON(&body); err != nil {
		_ = ctx.AbortWithError(400, err)
	}

	in := models.Student{
    FullName: body.FullName,
	}

	student, err := i.StudentRepository.Create(in)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	response := api.StudentsCreateResponse{
		Student: student.AsApiStudent(),
	}

	ctx.JSON(http.StatusCreated, response)
}

// StudentsGet implements the studentsGet contract from the OpenAPI spec.
func (i *OpenSchoolImpl) StudentsGet(ctx *gin.Context, id api.Cuid) {
	student, err := i.StudentRepository.Get(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if student == nil {
		_ = ctx.AbortWithError(http.StatusNotFound, errors.New("Student not found"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"student": student.AsApiStudent()})
}

func (i *OpenSchoolImpl) StudentsUpdate(ctx *gin.Context, id api.Cuid) {
	var body api.StudentsUpdateJSONRequestBody
	_ = ctx.Bind(&body)

	student := &models.Student{}
	student.Id = id
	student.FullName = body.FullName

	student, err := i.StudentRepository.Update(student)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, api.StudentsUpdateResponse{Student: student.AsApiStudent()})
}

func (i *OpenSchoolImpl) StudentsDelete(ctx *gin.Context, id api.Cuid) {
	student := models.Student{}
	student.Id = id

	err := i.StudentRepository.Delete(student)
	if err != nil {
		if err == students.StudentDoesNotExist {
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
