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

type SchemaType struct {
	typeName string
}

func (s SchemaType) String() string {
	return s.typeName
}

func (s *Schema) Type() SchemaType {
	return SchemaType{s.Name}
}

func (s *Schema) AddFields(fields ...Field) {
	s.Fields = append(s.Fields, fields...)
}

func (s *Schema) MutableColumns() []string {
	fields := []string{}
	for i := range s.Fields {
		if s.Fields[i].IsMutable() {
			fields = append(fields, gooostrings.ToSnakeCase(s.Fields[i].Name))
		}
	}

	return fields
}

func (s *Schema) ImmutableColumns() []string {
	fields := []string{}
	for i := range s.Fields {
		if s.Fields[i].IsImmutable() {
			fields = append(fields, gooostrings.ToSnakeCase(s.Fields[i].Name))
		}
	}

	return fields
}

func (s *Schema) SetClause() []string {
	placeholders := []string{}
	for i, c := range s.MutableColumns() {
		placeholders = append(placeholders, fmt.Sprintf("%s = $%d", gooostrings.ToSnakeCase(c), i+1))
	}

	for _, c := range s.ImmutableColumns() {
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
		if s.Fields[i].IsMutable() {
			placeholders = append(placeholders, fmt.Sprintf("$%d", index))
			index++
		}
	}

	return placeholders
}

func (s *Schema) ImmutablePlaceholders() []string {
	placeholders := []string{}
	index := 1
	for i := range s.Fields {
		if s.Fields[i].IsImmutable() {
			placeholders = append(placeholders, fmt.Sprintf("$%d", index))
			index++
		}
	}

	return placeholders
}

func (s *Schema) IgnoredFields() []Field {
	fields := []Field{}
	for i := range s.Fields {
		if s.Fields[i].Options.Ignore {
			fields = append(fields, s.Fields[i])
		}
	}

	return fields
}

func (s *Schema) ColumnFields() []Field {
	fields := []Field{}
	for i := range s.Fields {
		if !s.Fields[i].Options.Ignore {
			fields = append(fields, s.Fields[i])
		}
	}

	return fields
}

func (s *Schema) Columns() []string {
	fields := []string{}
	for _, f := range s.ColumnFields() {
		fields = append(fields, f.ColumnName())
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
		if s.Fields[i].IsMutable() {
			fields = append(fields, s.Fields[i])
		}
	}

	return fields
}

func (s *Schema) MutableFieldKeys() []string {
	fields := []string{}
	for i := range s.Fields {
		if s.Fields[i].IsMutable() {
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

func (f Field) String() string {
	str := ""
	field := fmt.Sprintf("\t%s %s", f.Name, f.Type)
	if f.Tag != "" {
		str = fmt.Sprintf("%s `%s`\n", field, f.Tag)
	} else {
		str = fmt.Sprintf("%s\n", field)
	}

	return str
}

func (f Field) ColumnName() string {
	return gooostrings.ToSnakeCase(f.Name)
}

func (f Field) IsMutable() bool {
	return !f.Options.Immutable && !f.Options.Ignore
}

func (f Field) IsImmutable() bool {
	return f.Options.Immutable && !f.Options.Ignore
}

type Validator struct {
	Fields   []string
	Validate validator.Validator
}

type FieldOptions struct {
	Immutable  bool
	PrimaryKey bool
	Ignore     bool // ignore fields for insert and update like fields of association.
	Validators []Validator
}
