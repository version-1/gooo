package response

import (
	"net/http"
)

type Adapter interface {
	Render(w http.ResponseWriter, payload any, status int) error
	Error(w http.ResponseWriter, err error, status int)
}

type Void struct{}

type Response[O any] struct {
	http.ResponseWriter
	status  int
	adapter Adapter
}

func New[O any](w http.ResponseWriter, a Adapter) *Response[O] {
	return &Response[O]{
		ResponseWriter: w,
		status:         http.StatusOK,
		adapter:        a,
	}
}

func (r Response[O]) Render(o O) {
	err := r.adapter.Render(r.ResponseWriter, o, r.status)
	if err != nil {
		r.adapter.Error(r.ResponseWriter, err, http.StatusInternalServerError)
	}
}

func (r *Response[O]) WriteHeader(code int) {
	r.ResponseWriter.WriteHeader(code)
	r.status = code
}

func (r Response[O]) renderError(err error) {
	r.adapter.Error(r.ResponseWriter, err, r.status)
}

func (r Response[O]) InternalServerError(err error) {
	r.status = http.StatusInternalServerError
	r.renderError(err)
}

func (r Response[O]) NotFound(err error) {
	r.status = http.StatusNotFound
	r.renderError(err)
}

func (r Response[O]) BadRequest(err error) {
	r.status = http.StatusBadRequest
	r.renderError(err)
}

func (r Response[O]) UnprocessableEntity(err error) {
	r.status = http.StatusUnprocessableEntity
	r.renderError(err)
}

func (r Response[O]) Unauthorized(err error) {
	r.status = http.StatusUnauthorized
	r.renderError(err)
}