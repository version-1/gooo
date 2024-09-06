package controller

import (
	"encoding/json"
	"net/http"
)

var _ http.ResponseWriter = &Response{}

type Response struct {
	ResponseWriter http.ResponseWriter
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
