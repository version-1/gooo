package adapter

import (
	"fmt"
	"net/http"

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

func (a JSONAPI) RenderError(w http.ResponseWriter, e any, options ...any) error {
	b, err := a.resolveError(e, options...)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	_, err = w.Write(b)
	return err
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

func (a JSONAPI) resolveError(e any, options ...any) ([]byte, error) {
	switch v := e.(type) {
	case jsonapi.Errors:
		s, err := jsonapi.NewErrors(v).Serialize()
		return []byte(s), err
	case jsonapi.Error:
		s, err := jsonapi.NewErrors(jsonapi.Errors{v}).Serialize()
		return []byte(s), err
	case jsonapi.ErrorCompatible:
		obj := v.ToJSONAPIError()
		s, err := jsonapi.NewErrors(jsonapi.Errors{obj}).Serialize()
		return []byte(s), err
	default:
		return []byte{}, JSONAPIInvalidTypeError{Payload: v}
	}
}
