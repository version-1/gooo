package adapter

import (
	"encoding/json"
	"net/http"
)

type Raw struct{}

func (a Raw) Render(w http.ResponseWriter, payload any, options ...any) error {
	return json.NewEncoder(w).Encode(payload)
}

func (a Raw) RenderError(w http.ResponseWriter, e error, options ...any) error {
	return json.NewEncoder(w).Encode(e)
}

func (a Raw) InternalServerError(w http.ResponseWriter, e error, options ...any) error {
	w.WriteHeader(http.StatusInternalServerError)
	return a.RenderError(w, e, options...)
}

func (a Raw) BadRequest(w http.ResponseWriter, e error, options ...any) error {
	w.WriteHeader(http.StatusBadRequest)
	return a.RenderError(w, e, options...)
}

func (a Raw) NotFound(w http.ResponseWriter, e error, options ...any) error {
	w.WriteHeader(http.StatusNotFound)
	return a.RenderError(w, e, options...)
}

func (a Raw) Unauthorized(w http.ResponseWriter, e error, options ...any) error {
	w.WriteHeader(http.StatusUnauthorized)
	return a.RenderError(w, e, options...)
}

func (a Raw) Forbidden(w http.ResponseWriter, e error, options ...any) error {
	w.WriteHeader(http.StatusForbidden)
	return a.RenderError(w, e, options...)
}
