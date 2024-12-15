package route

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/version-1/gooo/pkg/core/request"
	"github.com/version-1/gooo/pkg/core/response"
	"github.com/version-1/gooo/pkg/toolkit/middleware"
)

type HandlerFunc[I, O any] func(*response.Response[O], *request.Request[I])

var _ middleware.Handler = Handler[any, any]{}

type Handler[I, O any] struct {
	Path    string
	Method  string
	handler HandlerFunc[I, O]
	params  *Params
	adapter response.Adapter
}

func (h *Handler[I, O]) ShiftPath(base string) {
	h.Path = filepath.Clean(base + "/" + h.Path)
}

func (h Handler[I, O]) String() string {
	return fmt.Sprintf("[%s] %s", h.Method, h.Path)
}

func (h Handler[I, O]) Handler(res http.ResponseWriter, req *http.Request) {
	p := h.Param(*req.URL)
	customRequest := request.New[I](req, p)
	customResponse := response.New[O](res, h.adapter)
	h.handler(customResponse, customRequest)
}

func (h Handler[I, O]) Match(r *http.Request) bool {
	if r.Method != h.Method {
		return false
	}

	if r.URL.Path == h.Path {
		return true
	}

	parts := strings.Split(h.Path, "/")
	targetParts := strings.Split(r.URL.Path, "/")
	if len(parts) < len(targetParts) {
		return false
	}

	for i, part := range parts {
		if !strings.HasPrefix(part, ":") && part != targetParts[i] {
			return false
		}
	}

	return true
}

func (h Handler[I, O]) Param(uri url.URL) *Params {
	if h.params == nil {
		p := parseParams(h.Path, uri.Path)
		h.params = &p
	}
	return h.params
}
