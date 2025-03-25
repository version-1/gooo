package middleware

import (
	"fmt"
	"net/http"
)

type Middlewares []Middleware

func (m *Middlewares) Append(mw ...Middleware) {
	*m = append(*m, mw...)
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
	If   func(*http.Request) bool
	Do   func(http.ResponseWriter, *http.Request) bool
}

func (m Middleware) String() string {
	return fmt.Sprintf("Middleware %s", m.Name)
}

func Always(r *http.Request) bool {
	return true
}
