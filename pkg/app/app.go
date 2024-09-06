package app

import (
	gocontext "context"
	"net/http"
	"time"

	"github.com/version-1/gooo/pkg/config"
	"github.com/version-1/gooo/pkg/context"
	"github.com/version-1/gooo/pkg/controller"
	"github.com/version-1/gooo/pkg/http/response"
	"github.com/version-1/gooo/pkg/logger"
)

type Server struct {
	Addr         string
	Config       *config.App
	ErrorHandler func(w *response.Response, r *controller.Request, e error)
	Handlers     []controller.Handler
	Middlewares  []controller.Middleware
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
	rr := &controller.Request{
		Request: r,
	}
	ww := response.New(
		w,
		response.Options{
			Adapter: s.Config.DefaultResponseRenderer,
		},
	)

	for _, handler := range s.Handlers {
		if handler.Match(rr) {
			rr.Handler = handler
			return
		}
	}

	for _, m := range s.Middlewares {
		if m.If(rr) {
			s.withRecover("middleware", ww, rr, func() {
				if next := m.Do(ww, rr); !next {
					return
				}
			})
		}
	}

	s.withRecover("handler", ww, rr, func() {
		if rr.Handler.BeforeHandler != nil {
			(*rr.Handler.BeforeHandler)(ww, rr)
		}
		rr.Handler.Handler(ww, rr)
	})

	http.NotFound(w, r)
}

func WithDefaultMiddlewares(s *Server) {
	s.RegisterMiddlewares(
		controller.WithContext(
			func(r *controller.Request) *controller.Request {
				ctx := r.Context()
				ctx = context.WithAppConfig(ctx, s.Config)

				return r.WithContext(ctx)
			},
		),
		controller.RequestBodyLogger(s.Logger()),
		controller.RequestLogger(s.Logger()),
	)
}

func (s Server) Run(ctx gocontext.Context) {
	hs := &http.Server{
		Addr:           s.Addr,
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	defer hs.Shutdown(ctx)

	s.Logger().Infof("Server is running on %s", s.Addr)
	hs.ListenAndServe()
}

func (s Server) WalkThrough(cb func(h controller.Handler)) {
	for _, h := range s.Handlers {
		cb(h)
	}
}

func (s Server) withRecover(spot string, w *response.Response, r *controller.Request, fn func()) {
	defer func() {
		if e := recover(); e != nil {
			s.Logger().Errorf("Caught panic on %s", spot)
			if err, ok := e.(error); ok {
				s.ErrorHandler(w, r, err)
			}
		}
	}()

	fn()
}
