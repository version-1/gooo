package schema

import (
	"fmt"

	"github.com/version-1/gooo/pkg/datasource/orm/validator"
	gooostrings "github.com/version-1/gooo/pkg/strings"
)

type SchemaFactory struct {
	Primary       Field
	DefaultFields []Field
}

func (d SchemaFactory) NewSchema(fields []Field) *Schema {
	s := &Schema{}
	s.Fields = []Field{d.Primary}
	s.Fields = append(s.Fields, fields...)
	s.Fields = append(s.Fields, d.DefaultFields...)

	return s
}

type Schema struct {
	Name      string
	TableName string
	Fields    []Field
}

func (s *Schema) MutableColumns() []string {
	fields := []string{}
	for i := range s.Fields {
		if !s.Fields[i].Options.Immutable {
			fields = append(fields, gooostrings.ToSnakeCase(s.Fields[i].Name))
		}
	}

	return fields
}

func (s *Schema) ImmutableColumns() []string {
	fields := []string{}
	for i := range s.Fields {
		if s.Fields[i].Options.Immutable {
			fields = append(fields, gooostrings.ToSnakeCase(s.Fields[i].Name))
		}
	}

	return fields
}

func (s *Schema) SetClause() []string {
	placeholders := []string{}
	index := 1
	for i := range s.Fields {
		if !s.Fields[i].Options.Immutable {
			placeholders = append(placeholders, fmt.Sprintf("%s = $%d", gooostrings.ToSnakeCase(s.Fields[i].Name), index))
			index++
		}
	}

	for _, c := range s.Columns() {
		if c == "updated_at" {
			placeholders = append(placeholders, "updated_at = NOW()")
			return placeholders
		}
	}

	return placeholders
}

func (s *Schema) MutablePlaceholders() []string {
	placeholders := []string{}
	index := 1
	for i := range s.Fields {
		if !s.Fields[i].Options.Immutable {
			placeholders = append(placeholders, fmt.Sprintf("$%d", index))
			index++
		}
	}

	return placeholders
}

func (s *Schema) Columns() []string {
	fields := []string{}
	for i := range s.Fields {
		fields = append(fields, s.Fields[i].ColumnName())
	}

	return fields
}

func (s *Schema) FieldNames() []string {
	fields := []string{}
	for i := range s.Fields {
		fields = append(fields, s.Fields[i].Name)
	}

	return fields
}

func (s *Schema) MutableFields() []Field {
	fields := []Field{}
	for i := range s.Fields {
		if !s.Fields[i].Options.Immutable {
			fields = append(fields, s.Fields[i])
		}
	}

	return fields
}

func (s *Schema) MutableFieldKeys() []string {
	fields := []string{}
	for i := range s.Fields {
		if !s.Fields[i].Options.Immutable {
			fields = append(fields, gooostrings.ToSnakeCase(s.Fields[i].Name))
		}
	}

	return fields
}

func (s *Schema) PrimaryKey() string {
	for i := range s.Fields {
		if s.Fields[i].Options.PrimaryKey {
			return s.Fields[i].Name
		}
	}

	return ""
}

type Field struct {
	Name    string
	Type    FieldType
	Tag     string
	Options FieldOptions
}

func (f Field) ColumnName() string {
	return gooostrings.ToSnakeCase(f.Name)
}

type Validator struct {
	Fields   []string
	Validate validator.Validator
}

type FieldOptions struct {
	Immutable  bool
	PrimaryKey bool
	Validators []Validator
}
