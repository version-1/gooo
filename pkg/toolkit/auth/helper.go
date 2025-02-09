package auth

import (
	"net/http"

	"github.com/version-1/gooo/pkg/core/api/context"
)

func SetContextOnAuthorized[T any](r *http.Request, sub string, fetcher func(sub string) (T, error)) error {
	u, err := fetcher(sub)
	if err != nil {
		return err
	}

	r.WithContext(context.With(r.Context(), "user", u))
	return nil
}
