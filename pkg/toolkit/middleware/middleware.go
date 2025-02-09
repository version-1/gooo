package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/version-1/gooo/pkg/core/api/middleware"
	"github.com/version-1/gooo/pkg/toolkit/logger"
)

func RequestLogger(logger logger.Logger) middleware.Middleware {
	return middleware.Middleware{
		Name: "RequestLogger",
		If:   middleware.Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			logger.Infof("%s %s", r.Method, r.URL.Path)
			return true
		},
	}
}

func ResponseLogger(logger logger.Logger) middleware.Middleware {
	return middleware.Middleware{
		Name: "ResponseLogger",
		If:   middleware.Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			// FIXME: get stats code
			// logger.Infof("Status: %d", w.StatusCode())
			return true
		},
	}
}

func RequestBodyLogger(logger logger.Logger) middleware.Middleware {
	return middleware.Middleware{
		Name: "RequestBodyLogger",
		If: func(r *http.Request) bool {
			return r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch
		},
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
				logger.Errorf("Error reading request body: %v", err)
				return false
			}

			if len(b) > 0 {
				logger.Infof("body: %s", b)
			}

			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(b))

			return true
		},
	}
}

func RequestHeaderLogger(logger logger.Logger) middleware.Middleware {
	return middleware.Middleware{
		Name: "RequestHeaderLogger",
		If:   middleware.Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			logger.Infof("HTTP Headers: ")
			for k, v := range r.Header {
				logger.Infof("%s: %s", k, v)
			}
			return true
		},
	}
}

func CORS(origin, methods, headers []string) middleware.Middleware {
	return middleware.Middleware{
		Name: "CORS",
		If:   middleware.Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(origin, ", "))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))
			return true
		},
	}
}

func WithContext(callbacks ...func(r *http.Request) *http.Request) middleware.Middleware {
	return middleware.Middleware{
		Name: "WithContext",
		If:   middleware.Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			for _, cb := range callbacks {
				req := cb(r)
				*r = *req
			}

			return true
		},
	}
}

type Handler interface {
	fmt.Stringer
	Match(r *http.Request) bool
	Handler(w http.ResponseWriter, r *http.Request)
}

func RequestHandler(handlers []Handler) middleware.Middleware {
	return middleware.Middleware{
		Name: "RequestHandler",
		If:   middleware.Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			match := false
			for _, handler := range handlers {
				if handler.Match(r) {
					handler.Handler(w, r)
					match = true
					break
				}
			}
			if !match {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(fmt.Sprintf("Not found endpoint: %s", r.URL.Path)))
			}

			return match
		},
	}
}
