package app

import (
	"net/http"

	"github.com/version-1/gooo/pkg/core/api/context"
	"github.com/version-1/gooo/pkg/core/api/middleware"
	"github.com/version-1/gooo/pkg/core/api/route"

	helper "github.com/version-1/gooo/pkg/toolkit/middleware"
)

func WithDefaultMiddlewares(a *App, handlers ...route.HandlerInterface) middleware.Middlewares {
	_handlers := make([]helper.Handler, len(handlers))
	for i, h := range handlers {
		_handlers[i] = h
	}
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
		helper.RequestHandler(_handlers),
		helper.ResponseLogger(a.Logger()), // TODO: not implemented
	})

	return a.Middlewares
}
