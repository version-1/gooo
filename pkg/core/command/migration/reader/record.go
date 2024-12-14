package reader

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/version-1/gooo/pkg/command/migration/constants"
	"github.com/version-1/gooo/pkg/db"
)

type Record struct {
	Version       string
	Cache         string
	Kind          string
	RollbackQuery *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func IsUpSchema(r *Record) bool {
	return r.Kind == string(constants.UpMigration)
}

func scan(r *Record, rows *sql.Row) error {
	return rows.Scan(&r.Version, &r.Cache, &r.Kind, &r.RollbackQuery, &r.CreatedAt, &r.UpdatedAt)
}

func findRecord(ctx context.Context, db *db.DB, r *Record) error {
	query := fmt.Sprintf("SELECT version, cache, kind, rollback_query, created_at, updated_at FROM %s WHERE version = $1 ORDER BY version desc LIMIT 1", constants.ConfigTableName)
	if err := scan(r, db.QueryRowContext(ctx, query, r.Version)); err != nil {
		return err
	}

	return nil
}

func (r Record) Save(ctx context.Context, db *db.DB) error {
	existing := Record{
		Version: r.Version,
	}
	err := findRecord(ctx, db, &existing)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET cache = $1, kind = $3, rollback_query = $4 WHERE version = $2;", constants.ConfigTableName)
	if err == sql.ErrNoRows {
		query = fmt.Sprintf("INSERT INTO %s (cache, version, kind, rollback_query) VALUES ($1, $2, $3, $4);", constants.ConfigTableName)
	}

	_, err = db.ExecContext(ctx, query, r.Cache, r.Version, r.Kind, r.RollbackQuery)
	return err
}

func (r *Record) Delete(ctx context.Context, db *db.DB) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE version = $1;", constants.ConfigTableName)
	_, err := db.ExecContext(ctx, query, r.Version)
	return err
}
