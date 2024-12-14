package route

import (
	"net/http"

	"github.com/version-1/gooo/pkg/core/response"
)

func JSON[I, O any]() *Handler[I, O] {
	return &Handler[I, O]{
		adapter: response.JSONAdapter{},
	}
}

func HTML[I, O any]() *Handler[I, O] {
	return &Handler[I, O]{
		adapter: response.HTMLAdapter{},
	}
}

func (h *Handler[I, O]) Get(path string, handler HandlerFunc[I, O]) *Handler[I, O] {
	h.Path = path
	h.Method = http.MethodGet
	h.handler = handler

	return h
}

func (h *Handler[I, O]) Post(path string, handler HandlerFunc[I, O]) *Handler[I, O] {
	h.Path = path
	h.Method = http.MethodPost
	h.handler = handler

	return h
}

func (h *Handler[I, O]) Patch(path string, handler HandlerFunc[I, O]) *Handler[I, O] {
	h.Path = path
	h.Method = http.MethodPatch
	h.handler = handler

	return h
}

func (h *Handler[I, O]) Delete(path string, handler HandlerFunc[I, O]) *Handler[I, O] {
	h.Path = path
	h.Method = http.MethodDelete
	h.handler = handler

	return h
}

func Post[I, O any](path string, handler HandlerFunc[I, O]) *Handler[I, O] {
	return &Handler[I, O]{
		Path:    path,
		Method:  http.MethodPost,
		handler: handler,
	}
}

func Get[I, O any](path string, handler HandlerFunc[I, O]) *Handler[I, O] {
	return &Handler[I, O]{
		Path:    path,
		Method:  http.MethodGet,
		handler: handler,
	}
}

func Put[I, O any](path string, handler HandlerFunc[I, O]) *Handler[I, O] {
	return &Handler[I, O]{
		Path:    path,
		Method:  http.MethodPut,
		handler: handler,
	}
}

func Patch[I, O any](path string, handler HandlerFunc[I, O]) *Handler[I, O] {
	return &Handler[I, O]{
		Path:    path,
		Method:  http.MethodPatch,
		handler: handler,
	}
}

func Delete[I, O any](path string, handler HandlerFunc[I, O]) *Handler[I, O] {
	return &Handler[I, O]{
		Path:    path,
		Method:  http.MethodDelete,
		handler: handler,
	}
}
