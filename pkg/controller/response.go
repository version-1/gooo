package controller

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	w http.ResponseWriter
}

func (r *Response) JSON(payload json.Marshaler) *Response {
	json.NewEncoder(r.w).Encode(payload)

	return r
}

func (r *Response) Body(payload string) *Response {
	r.w.Write([]byte(payload))

	return r
}

func (r *Response) Status(code int) *Response {
	r.w.WriteHeader(code)
	return r
}

func (r *Response) Header(key, value string) *Response {
	r.w.Header().Set(key, value)
	return r
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
