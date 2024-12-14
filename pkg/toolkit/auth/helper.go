package auth

import (
	"github.com/version-1/gooo/pkg/context"
	"github.com/version-1/gooo/pkg/http/request"
)

func SetContextOnAuthorized[T any](r *request.Request, sub string, fetcher func(sub string) (T, error)) error {
	u, err := fetcher(sub)
	if err != nil {
		return err
	}

	r.WithContext(context.WithUserConfig(r.Context(), u))
	return nil
}
