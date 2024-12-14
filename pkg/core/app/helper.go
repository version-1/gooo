package app

import (
	"net/http"

	"github.com/version-1/gooo/pkg/core/context"
	"github.com/version-1/gooo/pkg/core/middleware"

	helper "github.com/version-1/gooo/pkg/toolkit/middleware"
)

func WithDefaultMiddlewares(a *App, handlers []helper.Handler) middleware.Middlewares {
	a.Middlewares = middleware.Middlewares([]middleware.Middleware{
		helper.WithContext(
			func(r *http.Request) *http.Request {
				ctx := r.Context()
				ctx = context.With(ctx, context.APP_CONFIG_KEY, a.Config)

				return r.WithContext(ctx)
			},
		),
		helper.RequestLogger(a.Logger()),
		helper.RequestBodyLogger(a.Logger()),
		helper.RequestHandler(handlers),
		helper.ResponseLogger(a.Logger()),
	})

	return a.Middlewares
}
