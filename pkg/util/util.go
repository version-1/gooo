package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/version-1/gooo/pkg/errors"
)

func LookupGomodDirPath() (string, error) {
	path, err := LookupFile("go.mod")
	return filepath.Dir(path), err
}

func LookupGomodPath() (string, error) {
	return LookupFile("go.mod")
}

func LookupFile(filename string) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err)
	}

	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		target := filepath.Join(strings.Join(parts[:i], "/"), filename)
		if _, err := os.Stat(target); err == nil {
			return target, nil
		}
	}

	return "", errors.Errorf("%s not found", filename)
}

func Basename(path string) string {
	l := filepath.Ext(path)
	return filepath.Base(path)[:len(filepath.Base(path))-len(l)]
}

func IsZero(v any) (bool, error) {
	switch vv := v.(type) {
	case string:
		return vv == "", nil
	case int:
		return vv == 0, nil
	case bool:
		return vv == false, nil
	case rune:
		return vv == 0, nil
	case uuid.UUID:
		return vv == uuid.Nil, nil
	default:
		return false, errors.Errorf("unsupported type: %T", v)
	}
}
