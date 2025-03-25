package jsonapi

import (
	"net/http"

	"github.com/google/uuid"
)

type CodeGetter interface {
	Code() string
}

type TitleGetter interface {
	Title() string
}

var _ Errable = ErrorResponse{}

type ErrorResponse struct {
	err error
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{err}
}

func (e ErrorResponse) Code() string {
	if c, ok := e.err.(CodeGetter); ok {
		return c.Code()
	}

	return "internal_server_error"
}

func (e ErrorResponse) Title() string {
	if c, ok := e.err.(TitleGetter); ok {
		return c.Title()
	}

	return "Internal Server Error"
}

func (e ErrorResponse) Error() string {
	return e.err.Error()
}

func (e ErrorResponse) ToJSONAPIError() Error {
	return Error{
		ID:     uuid.New().String(),
		Title:  e.Title(),
		Status: http.StatusInternalServerError,
		Code:   e.Code(),
		Detail: e.Error(),
	}
}

func NewInternalServerError(err error) Errable {
	return ErrorResponse{err}
}

func NewBadRequest(err error) Error {
	e := ErrorResponse{err}.ToJSONAPIError()
	e.Status = http.StatusBadRequest
	if _, ok := err.(CodeGetter); !ok {
		e.Code = "bad_request"
	}

	if _, ok := err.(TitleGetter); !ok {
		e.Title = "Bad Request"
	}

	return e
}

func NewUnauthorized(err error) Error {
	e := ErrorResponse{err}.ToJSONAPIError()
	e.Status = http.StatusUnauthorized
	if _, ok := err.(CodeGetter); !ok {
		e.Code = "unauthorized"
	}

	if _, ok := err.(TitleGetter); !ok {
		e.Title = "Unauthorized"
	}

	return e
}

func NewNotFound(err error) Error {
	e := ErrorResponse{err}.ToJSONAPIError()
	e.Status = http.StatusNotFound
	if _, ok := err.(CodeGetter); !ok {
		e.Code = "not_found"
	}

	if _, ok := err.(TitleGetter); !ok {
		e.Title = "Not Found"
	}

	return e
}

func NewForbidden(err error) Error {
	e := ErrorResponse{err}.ToJSONAPIError()
	e.Status = http.StatusForbidden
	if _, ok := err.(CodeGetter); !ok {
		e.Code = "forbidden"
	}

	if _, ok := err.(TitleGetter); !ok {
		e.Title = "Forbidden"
	}

	return e
}

type Resourcers []Resourcer

func (r Resourcers) ToJSONAPIResource() (Resources, Resources) {
	list := Resources{}
	includes := Resources{}
	for _, ele := range r {
		re, appending := ele.ToJSONAPIResource()
		list.Append(re)
		includes.Append(appending.Data...)
	}

	return list, includes
}
