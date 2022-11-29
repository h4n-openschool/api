package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
	"github.com/h4n-openschool/api/auth"
	"github.com/h4n-openschool/api/models"
	"golang.org/x/crypto/bcrypt"
)

func (i *OpenSchoolImpl) AuthCurrentUser(c *gin.Context) {
  i.Logger.Info("header is " + c.GetHeader("Authorization"))

  if ok := auth.MustAuthenticate(c, i.TeacherRepository); ok {
    return
  }
  teacher := c.Value("user").(*models.Teacher)

  c.JSON(200, teacher.AsApiTeacher())
}

func (i *OpenSchoolImpl) AuthLogin(c *gin.Context) {
	var body api.AuthLoginJSONRequestBody
	if err := c.BindJSON(&body); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}

	t, err := i.TeacherRepository.GetByEmail(body.Email)
	if err != nil {
    _ = c.AbortWithError(
      http.StatusInternalServerError,
      fmt.Errorf("failed to get teacher: %v", err.Error()),
    )
	}

	if t == nil {
    _ = c.AbortWithError(
      http.StatusInternalServerError,
      fmt.Errorf("teacher not found for email %v", body.Email),
    )
		_ = c.AbortWithError(http.StatusNotFound, errors.New("Teacher not found for that email."))
    return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(t.PasswordHash), []byte(body.Password)); err != nil {
    _ = c.AbortWithError(
      http.StatusInternalServerError,
      fmt.Errorf("failed to compare hash and password: %v", err.Error()),
    )
	}

	u := models.User[models.Teacher]{
		Person:   t,
		PersonId: t.Id,
	}
  token, err := u.Jwt()
  if err != nil {
    _ = c.AbortWithError(
      http.StatusInternalServerError,
      fmt.Errorf("failed to generate jwt: %v", err.Error()),
    )
  }

	c.JSON(http.StatusOK, api.AuthLoginResponse{
    Token: token,
  })
}
