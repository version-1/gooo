package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/version-1/gooo/pkg/logger"
)

type QueryRunner interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Tx interface {
	QueryRunner
	Commit() error
	Rollback() error
}

type DB struct {
	executor QueryRunner
	logger   QueryLogger
}

func New(conn QueryRunner) *DB {
	return &DB{executor: conn, logger: defaultLogger}
}

func (d *DB) SetLogger(l logger.Logger) {
	d.logger = &queryLoggerAdapter{l}
}

func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	d.logger.Log(query, args)
	return d.executor.QueryRow(query, args...)
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	d.logger.Log(query, args)
	return d.executor.Query(query, args...)
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.logger.Log(query, args)
	return d.executor.Exec(query, args...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	d.logger.Log(query, args)
	return d.executor.QueryContext(ctx, query, args...)
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	d.logger.Log(query, args)
	return d.executor.QueryRowContext(ctx, query, args...)
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	d.logger.Log(query, args)
	return d.executor.ExecContext(ctx, query, args...)
}

func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	id := uuid.New()
	switch v := d.executor.(type) {
	case *sql.DB:
		tx, err := v.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}

		d.logger.Println(fmt.Sprintf("begin tx: %s", id))
		return &txManager{id: id, DB: &DB{executor: tx, logger: d.logger}}, nil
	case *sqlx.DB:
		tx, err := v.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}

		d.logger.Println(fmt.Sprintf("begin tx: %s", id))
		return &txManager{id: id, DB: &DB{executor: tx, logger: d.logger}}, nil
	default:
		return nil, fmt.Errorf("executor doesn't implement BeginTx: %T", d.executor)
	}
}

type txManager struct {
	id uuid.UUID
	*DB
}

func (d *txManager) Commit() error {
	switch v := d.executor.(type) {
	case *sql.Tx:
		d.logger.Printf("commit: %s", d.id)
		return v.Commit()
	case *sqlx.Tx:
		d.logger.Printf("commit: %s", d.id)
		return v.Commit()
	default:
		return fmt.Errorf("executor doesn't implement Commit: %T", d.executor)
	}
}

func (d *txManager) Rollback() error {
	switch v := d.executor.(type) {
	case *sql.Tx:
		d.logger.Printf("rollback: %s", d.id)
		return v.Rollback()
	case *sqlx.Tx:
		d.logger.Printf("rollback: %s", d.id)
		return v.Rollback()
	default:
		return fmt.Errorf("executor doesn't implement Commit: %T", d.executor)
	}
}
