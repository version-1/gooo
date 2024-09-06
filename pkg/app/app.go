package app

import (
	"net/http"
	"time"

	"github.com/version-1/gooo/pkg/config"
	"github.com/version-1/gooo/pkg/context"
	"github.com/version-1/gooo/pkg/controller"
	"github.com/version-1/gooo/pkg/logger"
)

type Server struct {
	Addr        string
	Config      *config.App
	Handlers    []controller.Handler
	Middlewares []controller.Middleware
}

func (s *Server) SetLogger(l logger.Logger) {
	s.Config.Logger = l
}

func (s Server) Logger() logger.Logger {
	return s.Config.GetLogger()
}

func (s *Server) RegisterHandlers(h ...controller.Handler) {
	s.Handlers = append(s.Handlers, h...)
}

func (s *Server) RegisterMiddlewares(m ...controller.Middleware) {
	s.Middlewares = append(s.Middlewares, m...)
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, m := range s.Middlewares {
		if m.If(r) {
			if next := m.Do(w, r); !next {
				return
			}
		}
	}

	for _, handler := range s.Handlers {
		rr := &controller.Request{
			Request: r,
			Handler: handler,
		}
		if handler.Match(rr) {
			ww := &controller.Response{ResponseWriter: w}
			if handler.BeforeHandler != nil {
				(*handler.BeforeHandler)(ww, rr)
			}
			handler.Handler(ww, rr)
			return
		}
	}

	http.NotFound(w, r)
}

func WithDefaultMiddlewares(s *Server) {
	s.RegisterMiddlewares(
		controller.WithContext(
			func(r *http.Request) *http.Request {
				ctx := r.Context()
				ctx = context.WithAppConfig(ctx, s.Config)

				return r.WithContext(ctx)
			},
		),
		controller.RequestBodyLogger(s.Logger()),
		controller.RequestLogger(s.Logger()),
		controller.JSONResponse(),
	)
}

func (s Server) Run() {
	hs := &http.Server{
		Addr:           s.Addr,
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.Logger().Infof("Server is running on %s", s.Addr)
	hs.ListenAndServe()
}

func (s Server) WalkThrough(cb func(h controller.Handler)) {
	for _, h := range s.Handlers {
		cb(h)
	}
}
