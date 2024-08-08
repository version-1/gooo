package jsonapi

import (
	"fmt"
	"strings"
)

func isEmptyJSON(s string) bool {
	return s == "{}" || s == "[]" || s == "null"
}

type ObjectBuilder []Field

func (j ObjectBuilder) String() string {
	fields := []string{}
	for _, field := range j {
		fields = append(fields, field.String())
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(fields, ",\n"))
}

func (j ObjectBuilder) Append(fields ...Field) ObjectBuilder {
	j = append(j, fields...)

	return j
}

func NewField(key string, value any) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

type Field struct {
	Key   string
	Value any
}

func (f Field) String() string {
	return renderField(f.Key, f.Value)
}

func renderField(k string, v any) string {
	return fmt.Sprintf("\"%s\": %s", k, Stringify(v))
}

func renderList(v []string) string {
	return fmt.Sprintf("[\n%s\n]", strings.Join(v, ",\n"))
}
