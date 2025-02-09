package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/version-1/gooo/pkg/core/api/middleware"
)

type JWTAuth[T any] struct {
	If             func(r *http.Request) bool
	OnAuthorized   func(r *http.Request, sub string) error
	PrivateKey     *string
	TokenExpiresIn time.Duration
	Issuer         string
}

func (a JWTAuth[T]) GetPrivateKey() string {
	if a.PrivateKey == nil {
		return os.Getenv("GOOO_JWT_PRIVATE_KEY")
	}

	return *a.PrivateKey
}

func (a JWTAuth[T]) Sign(r *http.Request) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.TokenExpiresIn)),
		Issuer:    a.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.GetPrivateKey())
}

func (a JWTAuth[T]) Guard() middleware.Middleware {
	return middleware.Middleware{
		If: a.If,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			str := r.Header.Get("Authorization")
			token := strings.TrimSpace(strings.ReplaceAll(str, "Bearer ", ""))
			t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
				return a.GetPrivateKey(), nil
			})
			if err != nil {
				reportError(w, err)
				return false
			}

			expired, err := isExpired(t.Claims)
			if err != nil {
				reportError(w, err)
				return false
			}

			if expired {
				renderJSON(w, map[string]string{
					"code":   "auth:token_expired",
					"error":  "Unauthorized",
					"detail": err.Error(),
				}, http.StatusUnauthorized)
				return false
			}

			sub, err := t.Claims.GetSubject()
			if err != nil {
				reportError(w, err)
				return false
			}

			if err := a.OnAuthorized(r, sub); err != nil {
				reportError(w, err)
				return false
			}

			return true
		},
	}
}

func renderJSON(w http.ResponseWriter, payload map[string]string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func reportError(w http.ResponseWriter, e error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	payload := map[string]string{
		"code":   "unauthorized",
		"error":  "Unauthorized",
		"detail": e.Error(),
	}
	json.NewEncoder(w).Encode(payload)
}
