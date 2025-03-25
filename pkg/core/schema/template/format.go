package template

import (
	"fmt"
	"go/format"

	"golang.org/x/tools/imports"
)

func pretify(filename, s string) ([]byte, error) {
	formatted, err := format.Source([]byte(s))
	if err != nil {
		fmt.Println("Error processing format", s)
		return []byte{}, err
	}

	processed, err := imports.Process(filename, formatted, nil)
	if err != nil {
		fmt.Println("Error processing imports", s)
		return formatted, err
	}

	return processed, nil
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return string(s[0]-32) + s[1:]
}
