package controller

import (
	"encoding/json"
	"io"
	"net/http"
)

type Request struct {
	Handler Handler
	*http.Request
}

func MarshalBody[T json.Unmarshaler](r *Request, obj *T) error {
	b, err := io.ReadAll(r.Request.Body)
	if err != nil {
		return err
	}
	defer r.Request.Body.Close()

	return json.Unmarshal(b, obj)
}

func (r Request) Param(key string) (string, bool) {
	return r.Handler.Param(r.Request.URL.Path, key)
}

func (r Request) ParamInt(key string) (int, bool) {
	return r.Handler.ParamInt(r.Request.URL.Path, key)
}
