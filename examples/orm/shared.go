package orm

// this file is generated. DON'T EDIT this file
import (
	"context"
	"database/sql"
)

type scanner interface {
	Scan(dest ...any) error
}
type queryer interface {
	QueryRowContext(ctx context.Context, query string, dest ...any) *sql.Row
	QueryContext(ctx context.Context, query string, dest ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
