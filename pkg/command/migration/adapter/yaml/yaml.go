package yaml

import (
	"context"
	"fmt"

	"github.com/version-1/gooo/pkg/command/migration/constants"
	"github.com/version-1/gooo/pkg/command/migration/helper"
	"github.com/version-1/gooo/pkg/datasource/db"
)

type YamlElement interface {
	Load(path string) error
	Up(ctx context.Context, tx db.Tx) error
	Down(ctx context.Context, tx db.Tx) error
}

func LoadFile(path string) (*file, error) {
	f := &file{path: path}

	kind, err := f.Kind()
	if err != nil {
		return nil, err
	}

	switch kind {
	case constants.SchemaMigration:
		f.element = &OriginSchema{}
	case constants.UpMigration, constants.DownMigration:
		f.element = &RawSchema{}
	default:
		return nil, fmt.Errorf("invalid migration kind: %s", kind)
	}

	err = f.element.Load(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type file struct {
	path    string
	element YamlElement
}

func (f file) Path() string {
	return f.path
}

func (f file) Up(ctx context.Context, tx db.Tx) error {
	return f.element.Up(ctx, tx)
}

func (f file) Down(ctx context.Context, tx db.Tx) error {
	return f.element.Down(ctx, tx)
}

func (f file) Version() (string, error) {
	return helper.ParseVersion(f.path)
}

func (f file) Kind() (constants.MigrationKind, error) {
	return helper.ParseKind(f.path)
}

var _ YamlElement = (*OriginSchema)(nil)
var _ YamlElement = (*RawSchema)(nil)
