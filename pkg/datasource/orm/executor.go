package orm

import (
	"context"
	"database/sql"
)

var _ Tx = &Executor{}
var _ QueryRunner = &Executor{}

type Tx interface {
	QueryRunner
	Commit() error
	Rollback() error
}

type Executor struct {
	*Orm
	tx *sql.Tx
}

func NewExecutor(orm *Orm, tx ...*sql.Tx) *Executor {
	e := &Executor{Orm: orm}

	if len(tx) > 0 {
		e.tx = tx[0]
	}

	return e
}

func (e *Executor) queryRunner() QueryRunner {
	if e.tx != nil {
		return e.tx
	}

	return e.Orm
}

func (e *Executor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	e.Orm.LogQuery(query, args)

	return e.queryRunner().QueryRowContext(ctx, query, args...)
}

func (e *Executor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	e.Orm.LogQuery(query, args)

	return e.queryRunner().ExecContext(ctx, query, args...)
}

func (e *Executor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	e.Orm.LogQuery(query, args)

	return e.queryRunner().QueryContext(ctx, query, args...)
}

func (e *Executor) Commit() error {
	if e.tx == nil {
		return nil
	}

	return e.tx.Commit()
}

func (e *Executor) Rollback() error {
	if e.tx == nil {
		return nil
	}

	return e.tx.Rollback()
}
