package orm

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/version-1/gooo/pkg/datasource/logging"
)

type QueryRunner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Logger interface {
	Infof(string, ...any)
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Options struct {
	QueryLog bool
}

type Orm struct {
	db      *sqlx.DB
	logger  Logger
	ql      *logging.QueryLogger
	options Options
}

func New(db *sqlx.DB, logger Logger, options Options) *Orm {
	ql := logging.NewQueryLogger(logger)
	o := &Orm{
		db:      db,
		logger:  logger,
		ql:      ql,
		options: options,
	}

	return o
}

func (o *Orm) Transaction(ctx context.Context, fn func(*Executor) error) error {
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ex := NewExecutor(o, tx)
	if err = fn(ex); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}

func (o Orm) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	o.LogQuery(query, args)

	return o.db.QueryRowContext(ctx, query, args...)
}

func (o Orm) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	o.LogQuery(query, args)

	return o.db.ExecContext(ctx, query, args...)
}

func (o Orm) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	o.LogQuery(query, args)

	return o.db.QueryContext(ctx, query, args...)
}

func (o Orm) LogQuery(query string, args []any) {
	if !o.options.QueryLog {
		return
	}

	o.ql.Info(query, args...)
}

type Scanner interface {
	Scan(dest ...any) error
}
