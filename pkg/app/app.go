package app

import (
	gocontext "context"
	"net/http"
	"time"

	"github.com/version-1/gooo/pkg/config"
	"github.com/version-1/gooo/pkg/context"
	"github.com/version-1/gooo/pkg/controller"
	"github.com/version-1/gooo/pkg/http/request"
	"github.com/version-1/gooo/pkg/http/response"
	"github.com/version-1/gooo/pkg/logger"
)

type Server struct {
	Addr         string
	Config       *config.App
	ErrorHandler func(w *response.Response, r *request.Request, e error)
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
	rr := &request.Request{
		Request: r,
	}
	ww := response.New(
		w,
		response.Options{
			Adapter: s.Config.DefaultResponseRenderer,
		},
	)

	for _, m := range s.Middlewares {
		if m.If(rr) {
			s.withRecover(m.String(), ww, rr, func() {
				if next := m.Do(ww, rr); !next {
					return
				}
			})
		}
	}

	var target *controller.Handler
	for _, handler := range s.Handlers {
		if handler.Match(rr) {
			target = &handler
			break
		}
	}

	if target == nil {
		http.NotFound(w, r)
		return
	}

	rr.Handler = target
	s.withRecover(target.String(), ww, rr, func() {
		if target.BeforeHandler != nil {
			(*target.BeforeHandler)(ww, rr)
		}
		target.Handler(ww, rr)
	})
}

func WithDefaultMiddlewares(s *Server) {
	s.RegisterMiddlewares(
		controller.WithContext(
			func(r *request.Request) *request.Request {
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

func (s Server) withRecover(spot string, w *response.Response, r *request.Request, fn func()) {
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
