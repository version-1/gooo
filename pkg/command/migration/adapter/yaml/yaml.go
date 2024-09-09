package yaml

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
	"github.com/version-1/gooo/pkg/command/migration/constants"
	"github.com/version-1/gooo/pkg/command/migration/helper"
	"github.com/version-1/gooo/pkg/command/migration/reader"
	"github.com/version-1/gooo/pkg/db"
	"github.com/version-1/gooo/pkg/logger"
	yaml "gopkg.in/yaml.v3"
)

type SchemaManager struct {
	conn     *db.DB
	pathGlob string
	reader   *reader.SchemaReader
	elements YamlFiles
	logger   logger.Logger
}

func NewSchemaManager(conn *sqlx.DB, pathGlob string) (*SchemaManager, error) {
	d := db.New(conn)
	s := &SchemaManager{conn: d, pathGlob: pathGlob}
	s.reader = reader.New(d)
	if err := s.reader.Read(context.Background()); err != nil {
		return nil, err
	}

	matches, err := filepath.Glob(pathGlob)
	if err != nil {
		return nil, err
	}

	files := make(YamlFiles, len(matches))
	for i, m := range matches {
		f, err := LoadFile(m)
		if err != nil {
			return nil, err
		}

		files[i] = *f
	}

	sort.Sort(&files)
	s.elements = files

	s.logger = logger.DefaultLogger

	return s, nil
}

func (s *SchemaManager) SetLogger(l logger.Logger) {
	s.logger = l
}

func (s *SchemaManager) Up(ctx context.Context, tx db.Tx, size int) error {
	version, err := s.reader.Version()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	index := 0
	for _, f := range s.elements {
		if size > 0 && index >= size {
			break
		}

		fileVersion, err := f.Version()
		if err != nil {
			return err
		}

		k, err := f.Kind()
		if err != nil {
			return err
		}

		if k == constants.DownMigration {
			continue
		}

		if version != "" && fileVersion <= version {
			continue
		}

		s.logger.Infof("Applying migration: [UP]: %s", f.path)
		if err := f.element.Up(ctx, tx); err != nil {
			return err
		}

		r := &reader.Record{
			Version: fileVersion,
			Kind:    string(k),
		}

		if err := s.reader.Save(ctx, r); err != nil {
			return err
		}
		index++
	}

	return nil
}

func (s *SchemaManager) Down(ctx context.Context, tx db.Tx, size int) error {
	version, err := s.reader.Version()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	index := 0
	for i := range s.elements {
		if size > 0 && index >= size {
			break
		}

		// do migration in reverse order
		f := s.elements[len(s.elements)-1-i]
		fileVersion, err := f.Version()
		if err != nil {
			return err
		}

		k, err := f.Kind()
		if err != nil {
			return err
		}

		if k == constants.UpMigration {
			continue
		}

		if version != "" && fileVersion > version {
			continue
		}

		s.logger.Infof("Applying migration: [DOWN]: %s", f.path)
		if err := f.element.Down(ctx, tx); err != nil {
			return err
		}
		r := &reader.Record{
			Version: fileVersion,
		}

		if err := r.Delete(ctx, s.conn); err != nil {
			return err
		}
		index++
	}

	return nil
}

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

func (f file) Version() (string, error) {
	return helper.ParseVersion(f.path)
}

func (f file) Kind() (constants.MigrationKind, error) {
	return helper.ParseKind(f.path)
}

var _ YamlElement = (*OriginSchema)(nil)
var _ YamlElement = (*RawSchema)(nil)

type YamlFiles []file

func (s *YamlFiles) Len() int {
	return len(*s)
}

func (s *YamlFiles) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

func (s *YamlFiles) Less(i, j int) bool {
	a, _ := (*s)[i].Version()
	b, _ := (*s)[j].Version()
	return a < b
}

type RawSchema struct {
	Query string `yaml:"query"`
}

func (s *RawSchema) Load(path string) error {
	return load(path, s)
}

func (s *RawSchema) Up(ctx context.Context, tx db.Tx) error {
	if _, exec := tx.ExecContext(ctx, s.Query); exec != nil {
		return exec
	}

	return nil
}

func (s *RawSchema) Down(ctx context.Context, tx db.Tx) error {
	if _, exec := tx.ExecContext(ctx, s.Query); exec != nil {
		return exec
	}

	return nil
}

func load(path string, schema any) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(f, schema)
	if err != nil {
		return err
	}

	return nil
}
