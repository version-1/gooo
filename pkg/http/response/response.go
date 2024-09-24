package response

import (
	"encoding/json"
	"net/http"

	"github.com/version-1/gooo/pkg/http/response/adapter"
)

var _ http.ResponseWriter = &Response{}

var jsonapiAdapter Renderer = adapter.JSONAPI{}
var rawAdapter Renderer = adapter.Raw{}

type Renderer interface {
	Render(w http.ResponseWriter, payload any, options ...any) error
	RenderError(w http.ResponseWriter, err error, options ...any) error
	InternalServerError(w http.ResponseWriter, err error, options ...any) error
	NotFound(w http.ResponseWriter, err error, options ...any) error
	BadRequest(w http.ResponseWriter, err error, options ...any) error
	Unauthorized(w http.ResponseWriter, err error, options ...any) error
	Forbidden(w http.ResponseWriter, err error, options ...any) error
}

type Options struct {
	Adapter string
}

type Response struct {
	ResponseWriter http.ResponseWriter
	adapter        Renderer
	options        Options
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
	}
}

func (r Response) Adapter() Renderer {
	if r.adapter != nil {
		return r.adapter
	}

	return rawAdapter
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

func (r *Response) Status(code int) *Response {
	r.ResponseWriter.WriteHeader(code)
	return r
}

func (r *Response) Render(payload any, options ...any) error {
	return r.Adapter().Render(r.ResponseWriter, payload, options...)
}

func (r *Response) RenderError(payload error, options ...any) error {
	return r.Adapter().Render(r.ResponseWriter, payload, options...)
}

func (r *Response) SetHeader(key, value string) *Response {
	r.ResponseWriter.Header().Set(key, value)
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
}

func (r *Response) InternalServerError() *Response {
	return r.Status(http.StatusInternalServerError)
}

func (r *Response) NotFound() *Response {
	return r.Status(http.StatusNotFound)
}

func (r *Response) BadRequest() *Response {
	return r.Status(http.StatusBadRequest)
}

func (r *Response) Unauthorized() *Response {
	return r.Status(http.StatusUnauthorized)
}

func (r *Response) Forbidden() *Response {
	return r.Status(http.StatusForbidden)
}

func (r *Response) InternalServerErrorWith(e error, options ...any) error {
	return r.Adapter().InternalServerError(r.ResponseWriter, e, options...)
}

func (r *Response) NotFoundWith(e error, options ...any) error {
	return r.Adapter().NotFound(r.ResponseWriter, e, options...)
}

func (r *Response) BadRequestWith(e error, options ...any) error {
	return r.Adapter().BadRequest(r.ResponseWriter, e, options...)
}

func (r *Response) UnauthorizedWith(e error, options ...any) error {
	return r.Adapter().Unauthorized(r.ResponseWriter, e, options...)
}

func (r *Response) ForbiddenWith(e error, options ...any) error {
	return r.Adapter().Forbidden(r.ResponseWriter, e, options...)
}
