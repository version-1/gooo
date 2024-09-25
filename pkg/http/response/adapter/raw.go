package adapter

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Raw struct {
	w http.ResponseWriter
}

func (a Raw) ContentType() string {
	return "text/plain"
}

func (a Raw) Render(payload any, options ...any) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(payload)
	if err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (a Raw) RenderError(e error, options ...any) ([]byte, error) {
	return a.Render(e.Error(), options...)
}
