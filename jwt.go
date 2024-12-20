package gu

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func CreateToken(user any, expireTime ...time.Time) (string, error) {
	expire := time.Now().AddDate(0, 0, 1)

	if len(expireTime) > 0 {
		expire = expireTime[0]
	}

	claims := jwt.MapClaims{
		"user":        user,
		"expire_time": expire,
	}

	accessJwtKey := []byte(JWTSECRET)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(accessJwtKey)
}
