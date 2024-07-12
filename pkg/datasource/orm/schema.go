package orm

import (
	"reflect"

	"github.com/version-1/gooo/pkg/datasource/orm/errors"
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
	Fields []Field
}

func (s *Schema) FieldKeys() []string {
	fields := []string{}
	for i := range s.Fields {
		fields = append(fields, gooostrings.ToSnakeCase(s.Fields[i].Name))
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

func (s *Schema) Validate(data any) errors.ValidationError {
	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Struct {
		return errors.NewNotStructError(data)
	}

	for i := range s.Fields {
		for j := range s.Fields[i].Options.Validators {
			f := rv.FieldByName(s.Fields[i].Name)
			v := f.Interface()

			if err := s.Fields[i].Options.Validators[j](v); err != nil {
				return err
			}
		}
	}

	return nil
}

type Field struct {
	Name    string
	Options FieldOptions
}

type FieldOptions struct {
	Immutable  bool
	Validators []validator.ValidateFunc
}
