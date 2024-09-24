package adapter

import (
	"fmt"
	"net/http"

	goooerrors "github.com/version-1/gooo/pkg/errors"
	"github.com/version-1/gooo/pkg/presenter/jsonapi"
)

type JSONAPI struct {
	payload any
	meta    jsonapi.Serializer
}

type JSONAPIOption struct {
	Meta jsonapi.Serializer
}

type JSONAPIInvalidTypeError struct {
	Payload any
}

func (e JSONAPIInvalidTypeError) Error() string {
	return fmt.Sprintf("Invalid payload type. Payload must implement jsonapi.Resourcer. got: %T", e.Payload)
}

func (a JSONAPI) Render(w http.ResponseWriter, payload any, options ...any) error {
	b, err := a.resolve(payload, options...)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	_, err = w.Write(b)
	return err
}

func (a JSONAPI) RenderError(w http.ResponseWriter, e error, options ...any) error {
	b, errors, err := a.resolveError(e, options...)
	if err != nil {
		fmt.Println("error ==========", err)
		return err
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	if len(errors) > 0 {
		w.WriteHeader(errors[0].Status)
	}
	_, err = w.Write(b)
	return err
}

func (a JSONAPI) InternalServerError(w http.ResponseWriter, e error, options ...any) error {
	err := jsonapi.NewInternalServerError(e)
	return a.RenderError(w, err, options...)
}

func (a JSONAPI) BadRequest(w http.ResponseWriter, e error, options ...any) error {
	err := jsonapi.NewBadRequest(e)
	return a.RenderError(w, err, options...)
}

func (a JSONAPI) NotFound(w http.ResponseWriter, e error, options ...any) error {
	err := jsonapi.NewNotFound(e)
	return a.RenderError(w, err, options...)
}

func (a JSONAPI) Unauthorized(w http.ResponseWriter, e error, options ...any) error {
	err := jsonapi.NewUnauthorized(e)
	return a.RenderError(w, err, options...)
}

func (a JSONAPI) Forbidden(w http.ResponseWriter, e error, options ...any) error {
	err := jsonapi.NewForbidden(e)
	return a.RenderError(w, err, options...)
}

func (a JSONAPI) resolve(payload any, options ...any) ([]byte, error) {
	var meta jsonapi.Serializer
	for _, opt := range options {
		if t, ok := opt.(*JSONAPIOption); ok {
			meta = t.Meta
		}
	}

	switch v := payload.(type) {
	case jsonapi.Resourcer:
		data, includes := v.ToJSONAPIResource()
		s, err := jsonapi.New(data, includes, meta).Serialize()
		return []byte(s), err
	case []jsonapi.Resourcer:
		list := jsonapi.Resources{}
		includes := jsonapi.Resources{}
		for _, ele := range v {
			r, appending := ele.ToJSONAPIResource()
			list.Append(r)
			includes.Append(appending.Data...)
		}
		s, err := jsonapi.NewMany(list, includes, meta).Serialize()

		return []byte(s), err
	default:
		return []byte{}, JSONAPIInvalidTypeError{Payload: v}
	}
}

func (a JSONAPI) resolveError(e error, options ...any) ([]byte, []jsonapi.Error, error) {
	switch v := e.(type) {
	case jsonapi.Errors:
		s, err := jsonapi.NewErrors(v).Serialize()
		return []byte(s), v, err
	case jsonapi.Error:
		errors := jsonapi.Errors{v}
		s, err := jsonapi.NewErrors(errors).Serialize()
		return []byte(s), errors, err
	case jsonapi.Errable:
		obj := v.ToJSONAPIError()
		errors := jsonapi.Errors{obj}
		s, err := jsonapi.NewErrors(errors).Serialize()
		return []byte(s), errors, goooerrors.New(err.Error())
	default:
		obj := jsonapi.NewErrorResponse(v).ToJSONAPIError()
		errors := jsonapi.Errors{obj}
		s, err := jsonapi.NewErrors(errors).Serialize()
		return []byte(s), errors, err
	}
}
