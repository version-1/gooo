package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/version-1/gooo/pkg/http/request"
	"github.com/version-1/gooo/pkg/http/response"
	"github.com/version-1/gooo/pkg/logger"
)

type Middlewares []Middleware

func (m *Middlewares) Append(mw ...Middleware) {
	*m = append(*m, mw...)
}

func (m *Middlewares) Insert(index int, mw Middleware) {
	list := []Middleware{}
	for i, it := range *m {
		if i == index {
			list = append(list, mw)
		}

		list = append(list, it)
	}

	*m = list
}

func (m *Middlewares) Prepend(mw ...Middleware) {
	list := mw
	for _, it := range *m {
		list = append(list, it)
	}
	*m = list
}

type Middleware struct {
	Name string
	If   func(*request.Request) bool
	Do   func(*response.Response, *request.Request) bool
}

func (m Middleware) String() string {
	return fmt.Sprintf("Middleware %s", m.Name)
}

func Always(r *request.Request) bool {
	return true
}

func RequestLogger(logger logger.Logger) Middleware {
	return Middleware{
		If: Always,
		Do: func(w *response.Response, r *request.Request) bool {
			logger.Infof("%s %s", r.Request.Method, r.Request.URL.Path)
			return true
		},
	}
}

func ResponseLogger(logger logger.Logger) Middleware {
	return Middleware{
		If: Always,
		Do: func(w *response.Response, r *request.Request) bool {
			logger.Infof("Status: %d", w.StatusCode())
			return true
		},
	}
}

func RequestBodyLogger(logger logger.Logger) Middleware {
	return Middleware{
		If: Always,
		Do: func(w *response.Response, r *request.Request) bool {
			b, err := io.ReadAll(r.Request.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
				logger.Errorf("Error reading request body: %v", err)
				return false
			}

			io.Copy(w, io.MultiReader(bytes.NewReader(b), r.Request.Body))
			if len(b) > 0 {
				logger.Infof("body: %s", b)
			}
			return true
		},
	}
}

func RequestHeaderLogger(logger logger.Logger) Middleware {
	return Middleware{
		If: Always,
		Do: func(w *response.Response, r *request.Request) bool {
			logger.Infof("HTTP Headers: ")
			for k, v := range r.Request.Header {
				logger.Infof("%s: %s", k, v)
			}
			return true
		},
	}
}

func CORS(origin, methods, headers []string) Middleware {
	return Middleware{
		If: Always,
		Do: func(w *response.Response, r *request.Request) bool {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(origin, ", "))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))
			return true
		},
	}
}

func WithContext(callbacks ...func(r *request.Request) *request.Request) Middleware {
	return Middleware{
		If: Always,
		Do: func(w *response.Response, r *request.Request) bool {
			for _, cb := range callbacks {
				*r = *cb(r)
			}

			return true
		},
	}
}

func RequestHandler(handlers []Handler) Middleware {
	return Middleware{
		If: Always,
		Do: func(w *response.Response, r *request.Request) bool {
			match := false
			for _, handler := range handlers {
				if handler.Match(r) {
					if handler.BeforeHandler != nil {
						(*handler.BeforeHandler)(w, r)
					}
					handler.Handler(w, r)
					match = true
					break
				}
			}
			if !match {
				w.NotFoundWith(fmt.Errorf("Not found endpoint: %s", r.Request.URL.Path))
			}

			return match
		},
	}
}
