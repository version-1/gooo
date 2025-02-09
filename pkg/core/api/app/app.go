package app

import (
	gocontext "context"
	"net/http"
	"time"

	"github.com/version-1/gooo/pkg/core/api/middleware"
	"github.com/version-1/gooo/pkg/toolkit/errors"
	"github.com/version-1/gooo/pkg/toolkit/logger"
)

type App struct {
	Addr         string
	Config       *Config
	ErrorHandler func(w http.ResponseWriter, r *http.Request, e error)
	Middlewares  middleware.Middlewares
}

func (s *App) SetLogger(l logger.Logger) {
	s.Config.logger = l
}

func (s App) Logger() logger.Logger {
	return s.Config.Logger()
}

func (s App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, m := range s.Middlewares {
		if m.If(r) {
			s.withRecover(m.String(), w, r, func() {
				if next := m.Do(w, r); !next {
					return
				}
			})
		}
	}
}

func (s App) Run(ctx gocontext.Context) error {
	hs := &http.Server{
		Addr:           s.Addr,
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	defer hs.Shutdown(ctx)

	s.Logger().Infof("App is running on %s", s.Addr)
	return hs.ListenAndServe()
}

func (s App) withRecover(spot string, w http.ResponseWriter, r *http.Request, fn func()) {
	defer func() {
		if e := recover(); e != nil {
			s.Logger().Errorf("Caught panic on %s", spot)
			if err, ok := e.(error); ok {
				s.ErrorHandler(w, r, err)
			}

			if v, ok := e.(string); ok {
				err := errors.New(v)
				s.ErrorHandler(w, r, err)
			}
		}
	}()

	fn()
}
