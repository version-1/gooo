package helper

import (
	"fmt"
	"path/filepath"
	"strings"

	goooerrors "github.com/version-1/gooo/pkg/toolkit/errors"

	"github.com/version-1/gooo/pkg/command/migration/constants"
)

func ParseKind(path string) (constants.MigrationKind, error) {
	base := filepath.Base(path)
	parts := strings.Split(base, ".")
	if len(parts) < 3 {
		v, err := ParseVersion(path)
		if err != nil {
			return "", goooerrors.Wrap(err)
		}

		if v == strings.Repeat("0", 14) {
			return constants.SchemaMigration, nil
		}

		return "", fmt.Errorf("invalid migration kind: %s", path)
	}

	switch parts[1] {
	case "up":
		return constants.UpMigration, nil
	case "down":
		return constants.DownMigration, nil
	default:
		return "", fmt.Errorf("invalid migration kind: %s", parts[1])
	}
}

type InvalidVersionError struct {
	path string
}

func (e InvalidVersionError) Error() string {
	return fmt.Sprintf("invalid version: %s", e.path)
}

func ParseVersion(path string) (string, error) {
	base := filepath.Base(path)
	parts := strings.Split(base, "_")
	if len(parts) < 2 && parts[0] != strings.Repeat("0", 14) {
		return "", goooerrors.Wrap(InvalidVersionError{path})
	}

	return parts[0], nil
}
