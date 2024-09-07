package request

import (
	gocontext "context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/version-1/gooo/pkg/context"
	"github.com/version-1/gooo/pkg/logger"
)

type ParamParser interface {
	Param(url string, key string) (string, bool)
	ParamInt(url string, key string) (int, bool)
}

type Request struct {
	Handler ParamParser
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

func (r Request) Logger() logger.Logger {
	cfg := context.AppConfig(r.Request.Context())
	return cfg.Logger
}

func (r Request) Param(key string) (string, bool) {
	return r.Handler.Param(r.Request.URL.Path, key)
}

func (r Request) ParamInt(key string) (int, bool) {
	return r.Handler.ParamInt(r.Request.URL.Path, key)
}

func (r *Request) WithContext(ctx gocontext.Context) *Request {
	r.Request = r.Request.WithContext(ctx)
	return r
}