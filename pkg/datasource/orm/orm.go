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

type OrmFactory struct {
	*sqlx.DB
	logger  Logger
	ql      *logging.QueryLogger
	options Options
}

func NewOrmFactory(db *sqlx.DB, logger Logger, options Options) *OrmFactory {
	ql := logging.NewQueryLogger(logger)

	return &OrmFactory{
		DB:      db,
		logger:  logger,
		ql:      ql,
		options: options,
	}
}

func (f *OrmFactory) New(m Schema) *Orm {
	return &Orm{
		OrmFactory: f,
		m:          m,
	}
}

type Orm struct {
	*OrmFactory
	m Schema
}

//	func (o *Orm) Create(ctx context.Context, m Model) error {
//		return create(ctx, o, m)
//	}
//
//	func (o *Orm) Find(ctx context.Context, m Model) error {
//		return find(ctx, o, m)
//	}
//
//	func (o *Orm) Update(ctx context.Context, m Model) error {
//		return update(ctx, o, m)
//	}
//
//	func (o *Orm) Delete(ctx context.Context, m Model) error {
//		return delete(ctx, o, m)
//	}
func (o *Orm) Transaction(ctx context.Context, fn func(*Executor) error) error {
	tx, err := o.DB.BeginTx(ctx, nil)
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

func (o Orm) LogQuery(query string, args []any) {
	if !o.options.QueryLog {
		return
	}

	o.ql.Info(query, args...)
}

type Scanner interface {
	Scan(dest ...any) error
}
