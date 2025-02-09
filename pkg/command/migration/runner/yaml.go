package runner

import (
	"context"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
	"github.com/version-1/gooo/pkg/command/migration/adapter/yaml"
	"github.com/version-1/gooo/pkg/datasource/db"
)

type Yaml struct {
	runner   *Base
	pathGlob string
}

func NewYaml(pathGlob string) (*Yaml, error) {
	return &Yaml{
		pathGlob: pathGlob,
	}, nil
}

func (y *Yaml) Prepare(conn *sqlx.DB) error {
	r, err := New(db.New(conn))
	if err != nil {
		return err
	}

	matches, err := filepath.Glob(y.pathGlob)
	if err != nil {
		return err
	}

	files := make(Elements, len(matches))
	for i, m := range matches {
		f, err := yaml.LoadFile(m)
		if err != nil {
			return err
		}

		files[i] = *f
	}

	sort.Sort(&files)

	y.runner = r
	r.SetElements(files)
	return nil
}

func (y Yaml) Up(ctx context.Context, tx db.Tx, size int) error {
	return y.runner.Up(ctx, tx, size)
}

func (y Yaml) Down(ctx context.Context, tx db.Tx, size int) error {
	return y.runner.Down(ctx, tx, size)
}

func (y Yaml) BasePath() string {
	return filepath.Dir(y.pathGlob)
}

func (y Yaml) Ext() string {
	return "yaml"
}

func (y Yaml) Elements() Elements {
	return y.runner.Elements()
}
