package reader

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/version-1/gooo/pkg/command/migration/constants"
	"github.com/version-1/gooo/pkg/db"
	yaml "gopkg.in/yaml.v3"
)

type SchemaReader struct {
	db      *db.DB
	version string
	schema  *Schema
}

type Schema struct {
	Tables []Table `yaml:"tables" json:"tables"`
}

type Table struct {
	Name    string   `yaml:"name" json:"name"`
	Columns []Column `yaml:"columns" json:"columns"`
	Indexes []Index  `yaml:"indexes" json:"indexes"`
}

type Column struct {
	Name       string  `yaml:"name" json:"name"`
	Type       string  `yaml:"type" json:"type"`
	Default    *string `yaml:"default" json:"default"`
	Null       *bool   `yaml:"null" json:"null"`
	PrimaryKey *bool   `yaml:"primary_key" json:"primary_key"`
}

type Index struct {
	Name   string `yaml:"name" json:"name"`
	Column string `yaml:"column" json:"column"`
	Def    string `yaml:"def" json:"def"`
	Unique *bool  `yaml:"unique" json:"unique"`
	Pkey   *bool  `yaml:"is_pkey" json:"is_pkey"`
}

func New(conn *db.DB) *SchemaReader {
	return &SchemaReader{
		db: conn,
	}
}

const listTableQuery = `SELECT tablename FROM pg_catalog.pg_tables WHERE tablename NOT LIKE 'pg_%' and schemaname <> 'information_schema' AND schemaname <> 'gooo_migrations_meta'`
const listColumnsQuery = "SELECT column_name, udt_name, is_nullable, column_default FROM information_schema.columns WHERE table_name = $1;"
const listIndexesQuery = `SELECT
  i.indexrelid::regclass as index_name,
  ii.indexdef,
  i.indisprimary as pkey,
  i.indisunique as unique
FROM pg_index i
JOIN pg_class c on c.oid = i.indrelid
JOIN pg_class index_meta on index_meta.oid = i.indexrelid
JOIN pg_indexes ii on index_meta.relname = ii.indexname
WHERE c.relname = $1;`

func (r *SchemaReader) Read(ctx context.Context) error {
	rows, err := r.db.QueryContext(ctx, listTableQuery)
	if err != nil {
		return err
	}

	s := &Schema{}

	for rows.Next() {
		var table string
		if err = rows.Scan(&table); err != nil {
			return err
		}

		rows, err = r.db.QueryContext(ctx, listColumnsQuery, table)
		if err != nil {
			return err
		}

		t := Table{Name: table}
		for rows.Next() {
			c := Column{}
			isNullable := ""
			if err = rows.Scan(&c.Name, &c.Type, &isNullable, &c.Default); err != nil {
				return err
			}

			if isNullable == "YES" {
				null := true
				c.Null = &null
			} else if isNullable == "NO" {
				null := false
				c.Null = &null
			}

			t.Columns = append(t.Columns, c)
		}

		rows, err = r.db.QueryContext(ctx, listIndexesQuery, table)
		if err != nil {
			return err
		}
		for rows.Next() {
			i := Index{}
			if err = rows.Scan(
				&i.Name,
				&i.Def,
				&i.Pkey,
				&i.Unique,
			); err != nil {
				return err
			}
			t.Indexes = append(t.Indexes, i)
		}

		s.Tables = append(s.Tables, t)
	}

	r.schema = s

	return nil
}

func (r *SchemaReader) JSON() ([]byte, error) {
	if r.schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	return json.Marshal(r.schema)
}

func (r *SchemaReader) Yaml() ([]byte, error) {
	if r.schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	return yaml.Marshal(r.schema)
}

func (r *SchemaReader) latest(ctx context.Context, re *Record) error {
	query := fmt.Sprintf("SELECT version, cache, kind, rollback_query, created_at, updated_at FROM %s ORDER BY version desk LIMIT 1", constants.ConfigTableName)
	if err := scan(re, r.db.QueryRow(query)); err != nil {
		return err
	}

	return nil
}

func (r *SchemaReader) Schema(ctx context.Context) (*Schema, error) {
	if r.schema != nil {
		return r.schema, nil
	}

	re := Record{}
	err := r.latest(ctx, &re)
	if err != nil && err != sql.ErrNoRows {
		return r.schema, err
	}

	if err == sql.ErrNoRows {
		if err := r.Read(ctx); err != nil {
			return r.schema, err
		}

		b, err := r.JSON()
		if err != nil {
			return r.schema, err
		}
		re.Cache = string(b)
		re.Version = strings.Repeat("0", 14)

		if err := re.Save(ctx, r.db); err != nil {
			return r.schema, err
		}

		return r.schema, nil
	}

	r.schema = &Schema{}
	if err := json.Unmarshal([]byte(re.Cache), r.schema); err != nil {
		return r.schema, err
	}

	return r.schema, nil
}

func (r *SchemaReader) Save(ctx context.Context, re *Record) error {
	if err := r.Read(ctx); err != nil {
		return err
	}

	b, err := r.JSON()
	if err != nil {
		return err
	}
	re.Cache = string(b)
	// FIXME: calculate the diff between current and previous schema

	return re.Save(ctx, r.db)
}

func (r *SchemaReader) Version() (string, error) {
	if r.version != "" {
		return r.version, nil
	}

	query := fmt.Sprintf("SELECT version FROM %s ORDER BY version desc LIMIT 1", constants.ConfigTableName)
	if err := r.db.QueryRow(query).Scan(&r.version); err != nil {
		return "", err
	}

	return r.version, nil
}
