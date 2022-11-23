package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/models"
	"github.com/h4n-openschool/api/repos/teachers"
)

func MustAuthenticate(c *gin.Context, tr teachers.TeacherRepository) bool {
  var t *models.Teacher
  var err error

  tId := c.GetString("auth.userId")
  if tId == "" {
    _ = c.AbortWithError(401, errors.New("you are unauthenticated"))
    return true
  }

  if t, err = tr.Get(tId); err != nil {
    e := c.Error(err)
    _ = e.SetType(gin.ErrorTypePrivate)
    _ = c.AbortWithError(401, errors.New("you are unauthenticated"))
    return true
  }
  if t != nil {
    _ = c.AbortWithError(401, errors.New("you are unauthenticated"))
    return true
  }

  return false
}
