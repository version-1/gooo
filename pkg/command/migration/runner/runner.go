package runner

import (
	"context"
	"database/sql"

	"github.com/version-1/gooo/pkg/command/migration/constants"
	"github.com/version-1/gooo/pkg/command/migration/reader"
	"github.com/version-1/gooo/pkg/datasource/db"
	"github.com/version-1/gooo/pkg/toolkit/logger"
)

type Base struct {
	conn     *db.DB
	reader   *reader.SchemaReader
	elements Elements
	logger   logger.Logger
}

func New(conn *db.DB) (*Base, error) {
	ctx := context.Background()
	r := &Base{conn: conn}
	r.reader = reader.New(conn)
	if err := r.reader.Read(ctx); err != nil {
		return nil, err
	}

	r.logger = logger.DefaultLogger

	return r, nil
}

func (r *Base) SetLogger(l logger.Logger) {
	r.logger = l
}

func (r *Base) SetElements(elements Elements) {
	r.elements = elements
}

func (r Base) Elements() Elements {
	return r.elements
}

type Migration interface {
	Up(ctx context.Context, tx db.Tx) error
	Down(ctx context.Context, tx db.Tx) error
	Version() (string, error)
	Path() string
	Kind() (constants.MigrationKind, error)
}

type Elements []Migration

func (s *Elements) Len() int {
	return len(*s)
}

func (s *Elements) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

func (s *Elements) Less(i, j int) bool {
	a, _ := (*s)[i].Version()
	b, _ := (*s)[j].Version()
	return a < b
}

func (s *Base) Up(ctx context.Context, tx db.Tx, size int) error {
	version, err := s.reader.Version()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	index := 0
	for _, e := range s.elements {
		if size > 0 && index >= size {
			break
		}

		fileVersion, err := e.Version()
		if err != nil {
			return err
		}

		k, err := e.Kind()
		if err != nil {
			return err
		}

		if k == constants.DownMigration {
			continue
		}

		if version != "" && fileVersion <= version {
			continue
		}

		s.logger.Infof("Applying migration: [UP]: %s", e.Path())
		if err := e.Up(ctx, tx); err != nil {
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

func (r *Base) Down(ctx context.Context, tx db.Tx, size int) error {
	version, err := r.reader.Version()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	index := 0
	for i := range r.elements {
		if size > 0 && index >= size {
			break
		}

		// do migration in reverse order
		e := r.elements[len(r.elements)-1-i]
		fileVersion, err := e.Version()
		if err != nil {
			return err
		}

		k, err := e.Kind()
		if err != nil {
			return err
		}

		if k == constants.UpMigration {
			continue
		}

		if version != "" && fileVersion > version {
			continue
		}

		r.logger.Infof("Applying migration: [DOWN]: %s", e.Path())
		if err := e.Down(ctx, tx); err != nil {
			return err
		}
		re := &reader.Record{
			Version: fileVersion,
		}

		if err := re.Delete(ctx, r.conn); err != nil {
			return err
		}
		index++
	}

	return nil
}
