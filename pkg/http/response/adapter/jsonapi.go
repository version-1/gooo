package adapter

import (
	"fmt"

	goooerrors "github.com/version-1/gooo/pkg/errors"
	"github.com/version-1/gooo/pkg/presenter/jsonapi"
)

type JSONAPI struct {
	meta jsonapi.Serializer
}

type JSONAPIOption struct {
	Meta jsonapi.Serializer
}

type JSONAPIInvalidTypeError struct {
	Payload any
}

func (e JSONAPIInvalidTypeError) Error() string {
	return fmt.Sprintf("Payload must implement jsonapi.Resourcer or []jsonapi.Resourcer. got: %T", e.Payload)
}

func (a JSONAPI) ContentType() string {
	return "application/vnd.api+json"
}

func (a *JSONAPI) Render(payload any, options ...any) ([]byte, error) {
	return resolve(payload, options...)
}

func (a *JSONAPI) RenderError(e error, options ...any) ([]byte, error) {
	b, _, err := resolveError(e, options...)
	return b, err
}

func resolve(payload any, options ...any) ([]byte, error) {
	var meta jsonapi.Serializer
	for _, opt := range options {
		if t, ok := opt.(*JSONAPIOption); ok {
			meta = t.Meta
		}
	}

	switch v := payload.(type) {
	case jsonapi.Resourcerable:
		data, includes := v.Resourcer().ToJSONAPIResource()
		s, err := jsonapi.New(data, includes, meta).Serialize()

		return []byte(s), err
	case []jsonapi.Resourcerable:
		rl := []jsonapi.Resourcer{}
		for _, ele := range v {
			rl = append(rl, ele.Resourcer())
		}

		list, includes := jsonapi.Resourcers(rl).ToJSONAPIResource()
		s, err := jsonapi.NewMany(list, includes, meta).Serialize()

		return []byte(s), err
	case jsonapi.Resourcer:
		data, includes := v.ToJSONAPIResource()
		s, err := jsonapi.New(data, includes, meta).Serialize()
		return []byte(s), err
	case []jsonapi.Resourcer:
		list, includes := jsonapi.Resourcers(v).ToJSONAPIResource()
		s, err := jsonapi.NewMany(list, includes, meta).Serialize()

		return []byte(s), err
	case jsonapi.Resourcers:
		list, includes := v.ToJSONAPIResource()
		s, err := jsonapi.NewMany(list, includes, meta).Serialize()

		return []byte(s), err
	default:
		return []byte{}, goooerrors.Wrap(JSONAPIInvalidTypeError{Payload: v})
	}
}

func resolveError(e error, options ...any) ([]byte, []jsonapi.Error, error) {
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
		return []byte(s), errors, err
	default:
		obj := jsonapi.NewErrorResponse(v).ToJSONAPIError()
		errors := jsonapi.Errors{obj}
		s, err := jsonapi.NewErrors(errors).Serialize()
		return []byte(s), errors, err
	}
}
