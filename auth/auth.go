package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/h4n-openschool/api/repos/teachers"
)

func MustAuthenticate(c *gin.Context, tr teachers.TeacherRepository) bool {
	var err error

	tId := c.GetString("auth.userId")
	if tId == "" {
		_ = c.AbortWithError(401, errors.New("you are unauthenticated"))
		return true
	}

	t, err := tr.Get(tId)
	if err != nil {
		e := c.Error(err)
		_ = e.SetType(gin.ErrorTypePrivate)
		_ = c.AbortWithError(401, errors.New("you are unauthenticated"))
		return true
	}

	if t == nil {
		_ = c.AbortWithError(401, errors.New("you are unauthenticated"))
		return true
	}

	c.Set("user", t)

	return false
}
