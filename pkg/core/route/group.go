package route

import (
	"github.com/version-1/gooo/pkg/toolkit/middleware"
)

type GroupHandler struct {
	Path     string
	Handlers []HandlerInterface
}

type HandlerInterface interface {
	middleware.Handler
	ShiftPath(string)
}

func (g *GroupHandler) Add(h ...HandlerInterface) {
	g.Handlers = append(g.Handlers, h...)
}

func (g GroupHandler) List() []middleware.Handler {
	list := make([]middleware.Handler, len(g.Handlers))
	for i, h := range g.Handlers {
		h.ShiftPath(g.Path)
		list[i] = h
	}

	return list
}

func Walk(list []middleware.Handler, fn func(h middleware.Handler)) {
	for _, h := range list {
		fn(h)
	}
}
