package route

import (
	"github.com/version-1/gooo/pkg/toolkit/middleware"
)

type GroupHandler struct {
	Path string
	// We use HandlerInterface instead of route.Handler because route.Handler is generic,
	// which prevents us from determining the concrete type of the handler list.
	Handlers []HandlerInterface
}

type HandlerInterface interface {
	middleware.Handler
	ShiftPath(string) HandlerInterface
}

func (g *GroupHandler) Add(h ...HandlerInterface) {
	g.Handlers = append(g.Handlers, h...)
}

func (g GroupHandler) Children() []HandlerInterface {
	list := make([]HandlerInterface, len(g.Handlers))
	for i, h := range g.Handlers {
		shifted := h.ShiftPath(g.Path)
		list[i] = shifted
	}

	return list
}

func Walk(list []HandlerInterface, fn func(h middleware.Handler)) {
	for _, h := range list {
		fn(h)
	}
}
