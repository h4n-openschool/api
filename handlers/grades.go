package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
	"github.com/h4n-openschool/api/auth"
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/repos/grades"
	"github.com/h4n-openschool/api/utils"
)

// GradesList implements the gradesList operation from the OpenAPI
// specification in [../api/spec.yaml].
func (i *OpenSchoolImpl) GradesList(ctx *gin.Context, id api.Cuid, params api.GradesListParams) {
  if ok := auth.MustAuthenticate(ctx, i.TeacherRepository); ok {
    return
  }

	// Read pagination options from the GradesListParams object
	pagination := utils.NewPaginationQuery()
	pagination.ReadFromOptional(params.Page, params.PerPage)

	// Retrieve a paginated list of grades
	grades, err := i.GradeRepository.GetAll(id, pagination)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	// Get the total number of grades in the database to build the PaginationData object.
	total, err := i.GradeRepository.Count()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	// Generate pagination data from the total and input pagination options.
	paginationData := utils.GeneratePaginationData("/v1/grades", total, pagination)

	// Convert the grade model array to an api.GradeList type to meet the OpenAPI definition.
	gradeList := models.GradesAsApiGradeList(grades)

	// Build the response body as a GradesListResponse type.
	response := api.GradesListResponse{
		Grades:    gradeList,
		Pagination: paginationData,
	}

	// Respond to the request in JSON
	ctx.JSON(http.StatusOK, response)
}

// GradesCreate implements the gradesCreate contract from the OpenAPI spec.
func (i *OpenSchoolImpl) GradesCreate(ctx *gin.Context, classId api.Cuid) {
	var body api.GradesCreateRequest
	if err := ctx.BindJSON(&body); err != nil {
		_ = ctx.AbortWithError(400, err)
	}

	in := models.Grade{
    ClassId: classId,
    StudentId: body.StudentId,
    Value: body.Value,
	}

  grade, err := i.GradeRepository.Create(in)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	response := api.GradesCreateResponse{
		Grade: grade.AsApiGrade(),
	}

	ctx.JSON(http.StatusCreated, response)
}

// GradesGet implements the gradesGet contract from the OpenAPI spec.
func (i *OpenSchoolImpl) GradesGet(ctx *gin.Context, id api.Cuid, grade api.Cuid) {
	g, err := i.GradeRepository.Get(grade)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if g == nil {
		_ = ctx.AbortWithError(http.StatusNotFound, errors.New("grade not found"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"grade": g.AsApiGrade()})
}

func (i *OpenSchoolImpl) GradesUpdate(ctx *gin.Context, id api.Cuid, grade api.Cuid) {
	var body api.GradesUpdateJSONRequestBody
	_ = ctx.Bind(&body)

	g := &models.Grade{}
	g.Id = grade
  if body.Value != nil {
    g.Value = *body.Value
  }

	g, err := i.GradeRepository.Update(g)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, api.GradesUpdateResponse{Grade: g.AsApiGrade()})
}

func (i *OpenSchoolImpl) GradesDelete(ctx *gin.Context, id api.Cuid, grade api.Cuid) {
	g := models.Grade{}
	g.Id = grade

	err := i.GradeRepository.Delete(g)
	if err != nil {
		if err == grades.GradeDoesNotExist {
			_ = ctx.AbortWithError(http.StatusNotFound, err)
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})
}
