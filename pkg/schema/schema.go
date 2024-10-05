package schema

import (
	"fmt"

	"github.com/version-1/gooo/pkg/schema/internal/renderer"
	"github.com/version-1/gooo/pkg/schema/internal/valuetype"
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

func (s Schema) GetName() string {
	return s.Name
}

func (s Schema) GetTableName() string {
	return s.TableName
}

func (s *Schema) Type() SchemaType {
	return SchemaType{s.Name}
}

func (s *Schema) AddFields(fields ...Field) {
	s.Fields = append(s.Fields, fields...)
}

func (s Schema) MutableColumns() []string {
	fields := []string{}
	for i := range s.Fields {
		if s.Fields[i].IsMutable() {
			fields = append(fields, gooostrings.ToSnakeCase(s.Fields[i].Name))
		}
	}
	return fields
}

func (s Schema) MutableFieldNames() []string {
	fields := []string{}
	for i := range s.Fields {
		if s.Fields[i].IsMutable() {
			fields = append(fields, s.Fields[i].Name)
		}
	}
	return fields
}

func (s Schema) ImmutableColumns() []string {
	fields := []string{}
	for i := range s.Fields {
		if s.Fields[i].IsImmutable() {
			fields = append(fields, gooostrings.ToSnakeCase(s.Fields[i].Name))
		}
	}

	return fields
}

func (s Schema) SetClause() []string {
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
		if s.Fields[i].Tag.Ignore {
			fields = append(fields, s.Fields[i])
		}
	}

	return fields
}

func (s Schema) AtttributeFields() []Field {
	fields := []Field{}
	for i := range s.Fields {
		f := s.Fields[i]
		if !f.Tag.Ignore && !s.Fields[i].IsAssociation() && !f.Tag.PrimaryKey {
			fields = append(fields, s.Fields[i])
		}
	}

	return fields
}

func (s Schema) AttributeFieldNames() []string {
	fields := []string{}
	for i := range s.Fields {
		f := s.Fields[i]
		if !f.Tag.Ignore && !s.Fields[i].IsAssociation() && !f.Tag.PrimaryKey {
			fields = append(fields, s.Fields[i].Name)
		}
	}

	return fields
}

func (s Schema) FieldNames() []string {
	fields := []string{}
	for i := range s.Fields {
		fields = append(fields, s.Fields[i].Name)
	}

	return fields
}

func (s Schema) ColumnFields() []Field {
	fields := []Field{}
	for i := range s.Fields {
		f := s.Fields[i]
		if !f.Tag.Ignore && !s.Fields[i].IsAssociation() {
			fields = append(fields, s.Fields[i])
		}
	}

	return fields
}

func (s Schema) ColumnFieldNames() []string {
	fields := []string{}
	for i := range s.Fields {
		f := s.Fields[i]
		if !f.Tag.Ignore && !s.Fields[i].IsAssociation() {
			fields = append(fields, s.Fields[i].Name)
		}
	}

	return fields
}

func (s Schema) Columns() []string {
	fields := []string{}
	for _, f := range s.ColumnFields() {
		fields = append(fields, f.ColumnName())
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

func (s Schema) AssociationFields() []Field {
	fields := []Field{}
	for i := range s.Fields {
		if s.Fields[i].IsAssociation() {
			fields = append(fields, s.Fields[i])
		}
	}

	return fields
}

func (s Schema) AssociationFieldIdents() []renderer.AssociationIdent {
	idents := []renderer.AssociationIdent{}
	for i := range s.Fields {
		if s.Fields[i].IsAssociation() {
			field := s.Fields[i]
			t := fmt.Stringer(field.Type)
			ok := valuetype.MaySlice(t)
			if v, ok := t.(valuetype.Elementer); ok {
				t = v.Element()
			}

			typeName := gooostrings.ToSnakeCase(t.String())
			primaryKey := field.AssociationPrimaryKey()
			idents = append(idents, renderer.AssociationIdent{
				PrimaryKey:      primaryKey,
				FieldName:       field.Name,
				TypeName:        typeName,
				TypeElementExpr: field.TypeElementExpr,
				Slice:           ok,
				Ref:             field.IsRef(),
			})
		}
	}

	return idents
}

func (s Schema) PrimaryKey() string {
	for i := range s.Fields {
		if s.Fields[i].Tag.PrimaryKey {
			return s.Fields[i].Name
		}
	}

	return ""
}
