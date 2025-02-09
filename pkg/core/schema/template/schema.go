package template

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/version-1/gooo/pkg/core/schema/openapi/v3_0_0"
	"github.com/version-1/gooo/pkg/core/schema/openapi/yaml"
	"github.com/version-1/gooo/pkg/core/schema/template/partial"
	"github.com/version-1/gooo/pkg/toolkit/errors"
)

type SchemaFile struct {
	Schema      *v3_0_0.RootSchema
	PackageName string
	Content     string
}

func (s SchemaFile) Filename() string {
	return "internal/schema/schema"
}

func (s SchemaFile) Render() (string, error) {
	schemas := []Schema{}
	err := s.Schema.Components.Schemas.Each(func(key string, s v3_0_0.Schema) error {
		fields, err := extractFields(s.Properties, "")
		if err != nil {
			return err
		}

		schemas = append(schemas, Schema{
			Fields:   fields,
			TypeName: key,
		})
		return nil
	})
	if err != nil {
		return "", err
	}

	content, err := renderSchemas(schemas)
	if err != nil {
		return "", err
	}

	f := file{
		PackageName: s.PackageName,
		Content:     content,
	}

	tmpl := template.Must(template.New("file").ParseFS(tmpl, "components/file.go.tmpl"))
	var b bytes.Buffer
	if err := tmpl.ExecuteTemplate(&b, "file.go.tmpl", f); err != nil {
		return "", err
	}

	res, err := pretify(s.Filename(), b.String())
	if err != nil {
		return "", errors.Wrap(err)
	}
	return string(res), err
}

type Schema struct {
	Fields   []string
	TypeName string
}

func renderSchemas(schemas []Schema) (string, error) {
	var b bytes.Buffer
	for _, s := range schemas {
		tmpl := template.Must(template.New("struct").ParseFS(tmpl, "components/struct.go.tmpl"))
		if err := tmpl.ExecuteTemplate(&b, "struct.go.tmpl", s); err != nil {
			return "", errors.Wrap(err)
		}
		b.WriteString("\n")
	}
	return b.String(), nil
}

func extractFields(props yaml.OrderedMap[v3_0_0.Property], prefix string) ([]string, error) {
	var fields []string
	for i := 0; i < props.Len(); i++ {
		k, v := props.Index(i)
		key := formatKeyname(k)
		if v.Ref != "" {
			fields = append(fields, key+" "+pointer(schemaTypeName(v.Ref)))
			continue
		}

		t, err := extractFieldType(v, prefix)
		if err != nil {
			return []string{}, err
		}
		fields = append(fields, key+" "+t)
	}
	return fields, nil
}

func extractFieldType(prop v3_0_0.Property, prefix string) (string, error) {
	if prop.Ref != "" {
		return prefix + pointer(schemaTypeName(prop.Ref)), nil
	}

	switch {
	case isPrimitive(prop.Type):
		return prefix + convertGoType(prop.Type), nil
	case isDate(prop.Type):
		return prefix + "time.Time", nil
	case isObject(prop.Type):
		fields, err := extractFields(prop.Properties, prefix)
		if err != nil {
			return "", err
		}
		return prefix + partial.AnonymousStruct(fields), nil
	case isArray(prop.Type):
		if prop.Items == nil {
			return "", fmt.Errorf("Array must have items properties. %s\n", prop.Type)
		}
		return extractFieldType(*prop.Items, prefix+"[]")
	default:
		return "", fmt.Errorf("Unknown type: %s\n", prop.Type)
	}
}

func pointer(typeName string) string {
	return "*" + typeName
}

func formatKeyname(key string) string {
	if key == "id" {
		return strings.ToUpper(key)
	}

	if strings.HasSuffix(key, "Id") {
		return key[0:len(key)-2] + "ID"
	}

	return Capitalize(key)
}

func convertGoType(t string) string {
	m := map[string]string{
		"string":  "string",
		"number":  "int",
		"integer": "int",
		"boolean": "bool",
		"byte":    "[]byte",
	}
	v, ok := m[t]
	if !ok {
		return t
	}
	return v
}

func isPrimitive(t string) bool {
	primitives := map[string]bool{
		"string":  true,
		"number":  true,
		"integer": true,
		"boolean": true,
		"byte":    true,
	}
	_, ok := primitives[t]
	return ok
}

func isComplex(t string) bool {
	return !isPrimitive(t)
}

func isArray(t string) bool {
	return t == "array"
}

func isObject(t string) bool {
	return t == "object"
}

func isDate(t string) bool {
	return t == "date"
}
