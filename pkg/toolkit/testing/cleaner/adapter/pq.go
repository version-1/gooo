package adapter

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/version-1/gooo/pkg/datasource/db"
)

var excluded = pq.Array([]string{"schema_migrations"})

type Pq struct {
	conn db.Tx
}

func New(conn db.Tx) *Pq {
	return &Pq{conn: conn}
}

func (p Pq) ListTables(ctx context.Context) ([]string, error) {
	query := `SELECT
							table_name
						FROM
							information_schema.tables
						WHERE
							table_type = 'BASE TABLE'
							AND table_schema = 'public'
							AND table_name <> any($1)
							AND table_schema = 'public'
						ORDER BY table_name ASC
						`
	rows, err := p.conn.QueryContext(ctx, query, excluded)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	tables := []string{}
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return []string{}, err
		}

		tables = append(tables, t)
	}

	return tables, nil
}

func (p Pq) Truncate(ctx context.Context, table string) error {
	_, err := p.conn.ExecContext(ctx, "TRUNCATE TABLE "+table+" CASCADE")
	return err
}

func (p Pq) ResetIndexes(ctx context.Context, table string) error {
	index := table + "_pkey"

	if err := dropConstraint(ctx, p.conn, table, index); err != nil {
		panic(err)
	}

	q := "ALTER TABLE " + table + " ADD PRIMARY KEY (id)"
	if _, err := p.conn.ExecContext(ctx, q); err != nil {
		return err
	}

	if err := resetUniqueIndex(ctx, p.conn, table); err != nil {
		return err
	}

	return nil
}

func dropConstraint(ctx context.Context, tx db.Tx, table string, index string) error {
	q := "ALTER TABLE " + table + " DROP CONSTRAINT IF EXISTS " + index + " CASCADE"
	if _, err := tx.ExecContext(ctx, q); err != nil {
		return err
	}

	q = "DROP INDEX IF EXISTS " + index + " CASCADE"
	if _, err := tx.ExecContext(ctx, q); err != nil {
		return err
	}

	return nil
}

func resetUniqueIndex(ctx context.Context, tx db.Tx, table string) error {
	q := fmt.Sprintf(
		`
		  SELECT indexname, indexdef
				FROM pg_indexes
				WHERE tablename = '%s'
					AND indexdef LIKE '%%UNIQUE%%'
					AND indexname NOT LIKE '%%pkey%%'
				ORDER BY indexname ASC
		`,
		table,
	)
	rows, err := tx.QueryContext(ctx, q)
	if err != nil {
		return err
	}

	type obj struct {
		name string
		def  string
	}

	list := []obj{}
	for rows.Next() {
		var v obj
		if errr := rows.Scan(&v.name, &v.def); errr != nil {
			return errr
		}

		list = append(list, v)
	}
	defer rows.Close()

	for _, v := range list {
		if err := dropConstraint(ctx, tx, table, v.name); err != nil {
			return err
		}

		keys := getKeysFromIndexDef(v.def)
		query := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)", table, v.name, keys)
		if _, err := tx.ExecContext(ctx, query); err != nil {
			return err
		}
	}

	return nil
}

func getKeysFromIndexDef(str string) string {
	start := 0
	end := 0
	for i, s := range str {
		if string(s) == "(" {
			start = i + 1
		}

		if string(s) == ")" {
			end = i
		}
	}

	return str[start:end]
}
