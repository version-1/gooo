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
	RenderError(w http.ResponseWriter, err any, options ...any) error
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

func (r *Response) RenderError(payload any, options ...any) error {
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