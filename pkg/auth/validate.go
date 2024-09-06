package auth

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func isExpired(claims jwt.Claims) (bool, error) {
	date, err := claims.GetExpirationTime()
	if err != nil {
		return false, err
	}

	if date.After(time.Now()) {
		return true, NewTokenExpiredError(date.Time)
	}

	return false, nil
}
