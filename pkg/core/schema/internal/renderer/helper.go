package renderer

import (
	"fmt"
	"go/format"

	"github.com/version-1/gooo/pkg/errors"
	"golang.org/x/tools/imports"
)

func wrapQuote(list []string) []string {
	for i := range list {
		list[i] = fmt.Sprintf("\"%s\"", list[i])
	}

	return list
}

func pretify(filename, s string) (string, error) {
	// return s, nil
	formatted, err := format.Source([]byte(s))
	if err != nil {
		return s, errors.Wrap(err)
	}

	processed, err := imports.Process(filename, formatted, nil)
	if err != nil {
		return string(formatted), errors.Wrap(err)
	}

	return string(processed), nil
}
