package controller

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

type Middleware struct {
	If func(*http.Request) bool
	Do func(http.ResponseWriter, *http.Request) bool
}

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func Always(r *http.Request) bool {
	return true
}

func JSONResponse() Middleware {
	return Middleware{
		If: Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			w.Header().Set("Content-Type", "application/json")
			return true
		},
	}
}

func RequestLogger(logger Logger) Middleware {
	return Middleware{
		If: Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			logger.Infof("%s %s", r.Method, r.URL.Path)
			return true
		},
	}
}

func RequestBodyLogger(logger Logger) Middleware {
	return Middleware{
		If: Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
				logger.Errorf("Error reading request body: %v", err)
				return false
			}

			io.Copy(w, io.MultiReader(bytes.NewReader(b), r.Body))
			logger.Infof("body: %s", b)
			return true
		},
	}
}

func RequestHeaderLogger(logger Logger) Middleware {
	return Middleware{
		If: Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			for k, v := range r.Header {
				logger.Infof("%s: %s", k, v)
			}
			return true
		},
	}
}

func CORS(origin, methods, headers []string) Middleware {
	return Middleware{
		If: Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(origin, ", "))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))
			return true
		},
	}
}

func WithContext(callbacks ...func(r *http.Request) *http.Request) Middleware {
	return Middleware{
		If: Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			for _, cb := range callbacks {
				*r = *cb(r)
			}

			return true
		},
	}
}
