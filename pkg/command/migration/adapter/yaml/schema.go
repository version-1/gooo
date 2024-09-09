package yaml

import (
	"context"
	"fmt"
	"strings"

	"github.com/version-1/gooo/pkg/command/migration/constants"
	"github.com/version-1/gooo/pkg/db"
)

type OriginSchema struct {
	Tables []Table `yaml:"tables"`
}

type Table struct {
	Name    string   `yaml:"name"`
	Columns []Column `yaml:"columns"`
	Indexes []Index  `yaml:"indexes"`
}

func (t Table) Query() string {
	s := fmt.Sprintf("CREATE TABLE %s (", t.Name)
	for _, c := range t.Columns {
		s += c.Definition() + ", "
	}

	s = s[:len(s)-2] + ")"
	return s
}

type Column struct {
	Name       string  `yaml:"name"`
	Type       string  `yaml:"type"`
	Default    *string `yaml:"default"`
	Null       *bool   `yaml:"null"`
	PrimaryKey *bool   `yaml:"primary_key"`
}

func (c Column) Definition() string {
	s := fmt.Sprintf("%s %s", c.Name, c.Type)
	if c.Default != nil && (*c.Default) != "" {
		s += fmt.Sprintf(" DEFAULT %s", *c.Default)
	}

	if c.Null != nil && (*c.Null) == true {
		// do nothing
	} else {
		s += " NOT NULL"
	}

	if c.PrimaryKey != nil && (*c.PrimaryKey) {
		s += " PRIMARY KEY"
	}

	return s
}

type Index struct {
	Name       string      `yaml:"name"`
	Columns    []string    `yaml:"columns"`
	Unique     *bool       `yaml:"unique"`
	ForeignKey *ForeignKey `yaml:"foreign_key"`
}

type ForeignKey struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

func (i Index) Query(table string, kind constants.OperationKind) string {
	if i.ForeignKey != nil {
		return fmt.Sprintf(
			`ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s);`,
			table,
			i.Name,
			strings.Join(i.Columns, ", "),
			i.ForeignKey.Table,
			i.ForeignKey.Column,
		)
	}
	unique := ""
	if i.Unique != nil && (*i.Unique) {
		unique = "UNIQUE"
	}

	s := ""
	if kind == constants.AddOperationKind {
		s = fmt.Sprintf(
			"CREATE %s INDEX %s ON %s (%s)",
			unique,
			i.Name,
			table,
			strings.Join(i.Columns, ", "),
		)
	} else if kind == constants.DropOperationKind {
		s = fmt.Sprintf("DROP INDEX %s", i.Name)
	}

	return s
}

func (s *OriginSchema) Load(path string) error {
	return load(path, s)
}

func (s *OriginSchema) Up(ctx context.Context, db db.Tx) error {
	for _, t := range s.Tables {
		if _, err := db.ExecContext(ctx, t.Query()); err != nil {
			return err
		}

		for _, i := range t.Indexes {
			if _, err := db.ExecContext(
				ctx,
				i.Query(t.Name, constants.AddOperationKind),
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *OriginSchema) Down(ctx context.Context, db db.Tx) error {
	for _, t := range s.Tables {
		q := fmt.Sprintf("DROP TABLE %s CASCADE", t.Name)
		if _, err := db.ExecContext(ctx, q); err != nil {
			return err
		}
	}

	return nil
}
