package runner

import (
	"context"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
	"github.com/version-1/gooo/pkg/command/migration/adapter/yaml"
	"github.com/version-1/gooo/pkg/db"
)

type Yaml struct {
	runner   *Base
	pathGlob string
}

func NewYaml(conn *sqlx.DB, pathGlob string) (*Yaml, error) {
	r, err := New(db.New(conn))
	if err != nil {
		return nil, err
	}

	matches, err := filepath.Glob(pathGlob)
	if err != nil {
		return nil, err
	}

	files := make(Elements, len(matches))
	for i, m := range matches {
		f, err := yaml.LoadFile(m)
		if err != nil {
			return nil, err
		}

		files[i] = *f
	}

	sort.Sort(&files)
	r.SetElements(files)

	return &Yaml{
		runner:   r,
		pathGlob: pathGlob,
	}, nil
}

func (y Yaml) Up(ctx context.Context, tx db.Tx, size int) error {
	return y.runner.Up(ctx, tx, size)
}

func (y Yaml) Down(ctx context.Context, tx db.Tx, size int) error {
	return y.runner.Down(ctx, tx, size)
}
