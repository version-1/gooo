package response

import (
	"encoding/json"
	"net/http"

	"github.com/version-1/gooo/pkg/http/response/adapter"
	"github.com/version-1/gooo/pkg/logger"
)

var _ http.ResponseWriter = &Response{}

var jsonapiAdapter Renderer = &adapter.JSONAPI{}
var rawAdapter Renderer = &adapter.Raw{}

type Renderer interface {
	ContentType() string
	Render(payload any, options ...any) ([]byte, error)
	RenderError(err error, options ...any) ([]byte, error)
}

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type Options struct {
	Adapter string
	logger  Logger
}

type Response struct {
	ResponseWriter http.ResponseWriter
	adapter        Renderer
	options        Options
	statusCode     int
}

func New(r http.ResponseWriter, opts Options) *Response {
	adp := rawAdapter
	switch opts.Adapter {
	case "jsonapi":
		adp = jsonapiAdapter
	default:
		opts.Adapter = "raw"
	}

	return &Response{
		ResponseWriter: r,
		adapter:        adp,
		options:        opts,
		statusCode:     http.StatusOK,
	}
}

func (r Response) logger() Logger {
	if r.options.logger != nil {
		return r.options.logger
	}

	return logger.DefaultLogger
}

func (r *Response) Adapter() Renderer {
	if r.adapter == nil {
		r.adapter = rawAdapter
	}

	r.Header().Set("Content-Type", r.adapter.ContentType())
	return r.adapter
}

func (r *Response) SetAdapter(adp Renderer) *Response {
	r.adapter = adp
	return r
}

func (r *Response) JSON(payload any) *Response {
	r.Header().Set("Content-Type", "application/json")
	json.NewEncoder(r.ResponseWriter).Encode(payload)

	return r
}

func (r *Response) Body(payload string) *Response {
	r.ResponseWriter.Write([]byte(payload))

	return r
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) Render(payload any, options ...any) error {
	b, err := r.Adapter().Render(payload, options...)
	if err != nil {
		return err
	}

	_, err = r.Write(b)
	return err
}

func (r *Response) RenderError(payload error, options ...any) error {
	return r.renderErrorWith(func() {}, payload, options...)
}

func (r *Response) SetHeader(key, value string) *Response {
	r.Header().Set(key, value)
	return r
}

func (r Response) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *Response) Write(b []byte) (int, error) {
	return r.ResponseWriter.Write(b)
}

func (r *Response) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

func (r *Response) InternalServerError() {
	r.WriteHeader(http.StatusInternalServerError)
}

func (r *Response) NotFound() {
	r.WriteHeader(http.StatusNotFound)
}

func (r *Response) BadRequest() {
	r.WriteHeader(http.StatusBadRequest)
}

func (r *Response) Unauthorized() {
	r.WriteHeader(http.StatusUnauthorized)
}

func (r *Response) Forbidden() {
	r.WriteHeader(http.StatusForbidden)
}

func (r *Response) renderErrorWith(fn func(), e error, options ...any) error {
	r.logger().Errorf("%+v", e)
	b, err := r.Adapter().RenderError(e, options...)
	if err != nil {
		return err
	}

	fn()

	_, err = r.Write(b)
	return err
}

func (r *Response) InternalServerErrorWith(e error, options ...any) {
	err := r.renderErrorWith(r.InternalServerError, e, options...)
	if err != nil {
		r.logger().Errorf("got error on rendering internal_server_error")
		panic(err)
	}
}

func (r *Response) NotFoundWith(e error, options ...any) error {
	return r.renderErrorWith(r.NotFound, e, options...)
}

func (r *Response) BadRequestWith(e error, options ...any) error {
	return r.renderErrorWith(r.BadRequest, e, options...)
}

func (r *Response) UnauthorizedWith(e error, options ...any) error {
	return r.renderErrorWith(r.Unauthorized, e, options...)
}

func (r *Response) ForbiddenWith(e error, options ...any) error {
	return r.renderErrorWith(r.Forbidden, e, options...)
}
