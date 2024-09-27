package jsonapi

import (
	"encoding/json"

	goooerrors "github.com/version-1/gooo/pkg/errors"
)

func Stringify(v any) string {
	s, err := Escape(v)
	if err != nil {
		panic(err)
	}

	// Remove the quotes
	return s[1 : len(s)-1]
}

func Escape(i any) (string, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return "", goooerrors.Wrap(err)
	}

	return string(b), nil
}

func MustEscape(i any) string {
	s, err := Escape(i)
	if err != nil {
		panic(err)
	}

	return s
}
