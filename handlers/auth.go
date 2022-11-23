package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/api"
	"github.com/h4n-openschool/api/models"
	"golang.org/x/crypto/bcrypt"
)

func (i *OpenSchoolImpl) AuthLogin(c *gin.Context) {
	var body api.AuthLoginJSONRequestBody
	if err := c.BindJSON(&body); err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
	}

	t, err := i.TeacherRepository.GetByEmail(body.Email)
	if err != nil {
    _ = c.AbortWithError(
      http.StatusInternalServerError,
      errors.New(fmt.Sprintf("failed to get teacher: %v", err.Error())),
    )
	}
	if t == nil {
    _ = c.AbortWithError(
      http.StatusInternalServerError,
      errors.New(fmt.Sprintf("teacher not found for email %v", body.Email)),
    )
		_ = c.AbortWithError(http.StatusNotFound, errors.New("Teacher not found for that email."))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(t.PasswordHash), []byte(body.Password)); err != nil {
    _ = c.AbortWithError(
      http.StatusInternalServerError,
      errors.New(fmt.Sprintf("failed to compare hash and password: %v", err.Error())),
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
      errors.New(fmt.Sprintf("failed to generate jwt: %v", err.Error())),
    )
  }

	c.JSON(http.StatusOK, api.AuthLoginResponse{
    Token: token,
  })
}
