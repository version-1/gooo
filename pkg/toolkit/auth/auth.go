package auth

import (
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/version-1/gooo/pkg/controller"
	"github.com/version-1/gooo/pkg/http/request"
	"github.com/version-1/gooo/pkg/http/response"
)

type JWTAuth[T any] struct {
	If             func(r *request.Request) bool
	OnAuthorized   func(r *request.Request, sub string) error
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

func (a JWTAuth[T]) Sign(r *request.Request) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.TokenExpiresIn)),
		Issuer:    a.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.GetPrivateKey())
}

func (a JWTAuth[T]) Guard() controller.Middleware {
	return controller.Middleware{
		If: a.If,
		Do: func(w *response.Response, r *request.Request) bool {
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
				w.JSON(map[string]string{
					"code":   "auth:token_expired",
					"error":  "Unauthorized",
					"detail": err.Error(),
				})
				w.WriteHeader(http.StatusUnauthorized)
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

func reportError(w *response.Response, e error) {
	w.JSON(
		map[string]string{
			"code":   "unauthorized",
			"error":  "Unauthorized",
			"detail": e.Error(),
		},
	)
	w.WriteHeader(http.StatusUnauthorized)
}
