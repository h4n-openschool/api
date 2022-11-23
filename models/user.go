package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/h4n-openschool/api/utils"
)

type User[T interface{}] struct {
	BaseMetadata
	PersonId string
	Person   *T
}

func (u *User[T]) Jwt() (string, error) {
  claims := utils.UserClaims{}
  claims.RegisteredClaims = jwt.RegisteredClaims{
    Issuer:    `osapi`,
    Subject:   u.PersonId,
    Audience:  []string{"client", "server"},
    ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    IssuedAt:  jwt.NewNumericDate(time.Now()),
    NotBefore: jwt.NewNumericDate(time.Now()),
  }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(utils.JwtSigningKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
