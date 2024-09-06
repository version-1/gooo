package adapter

import (
	"encoding/json"
	"net/http"
)

type Raw struct{}

func (a Raw) Render(w http.ResponseWriter, payload any, options ...any) error {
	return json.NewEncoder(w).Encode(payload)
}

func (a Raw) RenderError(w http.ResponseWriter, e any, options ...any) error {
	return json.NewEncoder(w).Encode(e)
}
