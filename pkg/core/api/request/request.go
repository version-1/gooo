package request

import (
	gocontext "context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/version-1/gooo/pkg/core/api/context"
	"github.com/version-1/gooo/pkg/toolkit/logger"
)

type Void struct{}

type Params interface {
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	GetBool(key string) (bool, error)
}

type Request[I any] struct {
	params Params
	*http.Request
	body  *[]byte
	query Query
}

func New[I any](r *http.Request, p Params) *Request[I] {
	return &Request[I]{
		Request: r,
	}
}

func (r *Request[I]) Body() (I, error) {
	var res I
	if r.body == nil {
		b, err := io.ReadAll(r.Request.Body)
		if err != nil {
			r.Logger().Errorf("failed to read request body: %s", err)
			return res, err
		}

		r.body = &b
	}

	if err := json.Unmarshal(*r.body, &res); err != nil {
		r.Logger().Errorf("failed to unmarshal request body: %s", err)
		return res, nil
	}

	return res, nil
}

type loggerGetter interface {
	Logger() logger.Logger
}

func (r Request[I]) Logger() logger.Logger {
	cfg := context.Get[loggerGetter](r.Request.Context(), context.APP_CONFIG_KEY)
	return cfg.Logger()
}

func (r Request[I]) Params() Params {
	return r.params
}

func (r Request[I]) Query() Query {
	return r.query
}

func (r *Request[I]) WithContext(ctx gocontext.Context) *Request[I] {
	r.Request = r.Request.WithContext(ctx)
	return r
}
