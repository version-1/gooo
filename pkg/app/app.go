package app

import (
	"net/http"
	"time"

	"github.com/version-1/gooo/pkg/controller"
)

type Server struct {
	Addr        string
	handlers    []controller.Handler
	middlewares []controller.Middleware
}

func (s *Server) Register(h controller.Handler) {
	s.handlers = append(s.handlers, h)
}

func (s *Server) RegisterMiddleware(m controller.Middleware) {
	s.middlewares = append(s.middlewares, m)
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, m := range s.middlewares {
		if m.If(r) {
			if next := m.Do(w, r); !next {
				return
			}
		}
	}

	for _, handler := range s.handlers {
		rr := &controller.Request{
			Request: r,
			Handler: handler,
		}
		if handler.Match(rr) {
			if handler.BeforeHandler != nil {
				(*handler.BeforeHandler)(w, rr)
			}
			handler.Handler(w, rr)
			return
		}
	}

	http.NotFound(w, r)
}

func (s Server) Run() {
	hs := &http.Server{
		Addr:           s.Addr,
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	hs.ListenAndServe()
}
