package controller

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/version-1/gooo/pkg/http/response"
)

type BeforeHandlerFunc func(*response.Response, *Request) bool
type HandlerFunc func(*response.Response, *Request)

type GroupHandler struct {
	Path     string
	Handlers []Handler
}

func (g *GroupHandler) Add(h ...Handler) {
	g.Handlers = append(g.Handlers, h...)
}

func (g GroupHandler) List() []Handler {
	list := make([]Handler, len(g.Handlers))
	for i, h := range g.Handlers {
		h.Path = filepath.Clean(g.Path + h.Path)
		list[i] = h
	}

	return list
}

type Handler struct {
	Path          string
	Method        string
	BeforeHandler *BeforeHandlerFunc
	Handler       HandlerFunc
}

func (h Handler) Match(r *Request) bool {
	if r.Request.Method != h.Method {
		return false
	}

	if r.Request.URL.Path == h.Path {
		return true
	}

	parts := strings.Split(h.Path, "/")
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
	if !strings.Contains(h.Path, search) {
		return "", false
	}

	parts := strings.Split(h.Path, "/")
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
		Path:    path,
		Method:  "POST",
		Handler: handler,
	}
}

func Get(path string, handler HandlerFunc) Handler {
	return Handler{
		Path:    path,
		Method:  "GET",
		Handler: handler,
	}
}

func Put(path string, handler HandlerFunc) Handler {
	return Handler{
		Path:    path,
		Method:  "PUT",
		Handler: handler,
	}
}

func Patch(path string, handler HandlerFunc) Handler {
	return Handler{
		Path:    path,
		Method:  "PATCH",
		Handler: handler,
	}
}

func Delete(path string, handler HandlerFunc) Handler {
	return Handler{
		Path:    path,
		Method:  "DELETE",
		Handler: handler,
	}
}
