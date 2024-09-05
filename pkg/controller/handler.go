package controller

import (
	"net/http"
	"strconv"
	"strings"
)

type BeforeHandlerFunc func(http.ResponseWriter, *Request) bool
type HandlerFunc func(http.ResponseWriter, *Request)

type Handler struct {
	path          string
	method        string
	BeforeHandler *BeforeHandlerFunc
	Handler       HandlerFunc
}

func (h Handler) Match(r *Request) bool {
	if r.Request.Method != h.method {
		return false
	}

	if r.Request.URL.Path == h.path {
		return true
	}

	parts := strings.Split(h.path, "/")
	targetParts := strings.Split(r.Request.URL.Path, "/")
	if len(parts) < len(targetParts) {
		return false
	}

	for i, part := range parts {
		if !strings.HasPrefix(part, ":") && part != targetParts[i] {
			return false
		}
	}

	return false
}

func (h Handler) Param(url string, key string) (string, bool) {
	search := ":" + key
	if !strings.Contains(h.path, search) {
		return "", false
	}

	parts := strings.Split(h.path, "/")
	index := -1
	for i, part := range parts {
		if part == search {
			index = i
			break
		}
	}

	if index == -1 {
		return "", false
	}

	targetParts := strings.Split(url, "/")
	if len(targetParts) < index {
		return "", false
	}

	return targetParts[index], true
}

func (h Handler) ParamInt(url string, key string) (int, bool) {
	v, ok := h.Param(url, key)
	if !ok {
		return 0, false
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}

	return n, true
}

func Post(path string, handler HandlerFunc) Handler {
	return Handler{
		path:    path,
		method:  "POST",
		Handler: handler,
	}
}

func Get(path string, handler HandlerFunc) Handler {
	return Handler{
		path:    path,
		method:  "GET",
		Handler: handler,
	}
}

func Put(path string, handler HandlerFunc) Handler {
	return Handler{
		path:    path,
		method:  "PUT",
		Handler: handler,
	}
}

func Patch(path string, handler HandlerFunc) Handler {
	return Handler{
		path:    path,
		method:  "PATCH",
		Handler: handler,
	}
}

func Delete(path string, handler HandlerFunc) Handler {
	return Handler{
		path:    path,
		method:  "DELETE",
		Handler: handler,
	}
}
