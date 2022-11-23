package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
  JwtSigningKey = "lol u thought this would be secure"
)

type UserClaims struct {
	jwt.RegisteredClaims
}

func AuthenticateMiddleware(c *gin.Context) {
  token := c.GetHeader("Authorization")

  if token != "" {
    parts := strings.Split(token, " ")
    if parts[0] != "Bearer" {
      _ = c.AbortWithError(401, errors.New("Invalid token type"))
    }

    claims := UserClaims{}
    t, err := jwt.ParseWithClaims(parts[1], &claims, func (t *jwt.Token) (interface{}, error) {
      return []byte(JwtSigningKey), nil
    })

    if err != nil {
      _ = c.Error(gin.Error{Err: err, Type: gin.ErrorTypePrivate})
      _ = c.AbortWithError(401, errors.New("you are unauthenticated"))
    }

    c.Set("auth.token", t)
    c.Set("auth.userId", claims.Subject)
  }

  c.Next()
}
