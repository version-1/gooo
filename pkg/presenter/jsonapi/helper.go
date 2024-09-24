package jsonapi

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Attributes[T any] struct {
	v T
}

func NewAttributes[T any](v T) Attributes[T] {
	return Attributes[T]{v}
}

func (v Attributes[T]) JSONAPISerialize() (string, error) {
	if b, err := json.Marshal(v); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

type Resourcable interface {
	ID() string
	Type() string
	Resources() Resources
}

type ResourceTemplate struct {
	Target Resourcable
}

func (v ResourceTemplate) ToJSONAPIResource() (data Resource, included Resources) {
	t := v.Target
	return Resource{
		ID:         t.ID(),
		Type:       t.Type(),
		Attributes: NewAttributes(v),
	}, t.Resources()
}

type resourcerer interface {
	Resourcer() Resourcer
}

func ToResourcerList[T resourcerer](list []T) []Resourcer {
	resources := make([]Resourcer, 0, len(list))
	for _, r := range list {
		resources = append(resources, r.Resourcer())
	}
	return resources
}

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
		Title:  "Internal Server Error",
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
