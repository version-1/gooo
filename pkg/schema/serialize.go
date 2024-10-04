package schema

import (
	"fmt"
	"strings"

	"github.com/version-1/gooo/pkg/schema/internal/template"
	gooostrings "github.com/version-1/gooo/pkg/strings"
)

func (s SchemaTemplate) defineToJSONAPIResource() string {
	primaryKey := s.Schema.PrimaryKey()

	str := fmt.Sprintf(`includes := &jsonapi.Resources{ShouldSort: true}
		r := &jsonapi.Resource{
		  ID:   jsonapi.Stringify(obj.%s),
			Type: "%s",
			Attributes: obj,
			Relationships: jsonapi.Relationships{},
		}
	`, primaryKey, gooostrings.ToSnakeCase(s.Schema.Name))
	str += "\n"

	for _, field := range s.Schema.AssociationFields() {
		t := fmt.Stringer(field.Type)
		_, ok := t.(slice)
		if v, ok := t.(Elementer); ok {
			t = v.Element()
		}
		typeName := gooostrings.ToSnakeCase(t.String())
		primaryKey := field.AssociationPrimaryKey()
		if ok {
			str += fmt.Sprintf(
				`elements := []jsonapi.Resourcer{}
				for _, ele := range obj.%s {
					elements = append(elements, jsonapi.Resourcer(ele))
				}
				jsonapi.HasMany(r, includes, elements, "%s", func(ri *jsonapi.ResourceIdentifier, i int) {
						id := obj.%s[i].%s
						ri.ID = jsonapi.Stringify(id)
				})`,
				field.Name,
				typeName,
				field.Name,
				primaryKey,
			)
			str += "\n"
		} else {
			if field.IsRef() {
				str += fmt.Sprintf(
					`ele := obj.%s
					if ele != nil {
						jsonapi.HasOne(r, includes, ele, ele.%s, "%s")
					}`,
					field.Name,
					primaryKey,
					typeName,
				)
			} else {
				str += fmt.Sprintf(
					`ele := obj.%s
					if ele.%s != (%s{}).%s {
						jsonapi.HasOne(r, includes, ele, ele.%s, "%s")
					}`,
					field.Name,
					primaryKey,
					field.TypeElementExpr,
					primaryKey,
					primaryKey,
					typeName,
				)
			}
			str += "\n"
		}
		str += "\n"
	}

	str += "\n"
	str += "return *r, *includes"

	return template.Method{
		Receiver:    s.Schema.Name,
		Name:        "ToJSONAPIResource",
		Args:        []template.Arg{},
		ReturnTypes: []string{"jsonapi.Resource", "jsonapi.Resources"},
		Body:        str,
	}.String()
}

func (s SchemaTemplate) defineJSONAPISerialize() string {
	fields := []string{}
	for _, field := range s.Schema.AtttributeFields() {
		v := fmt.Sprintf(
			`fmt.Sprintf("\"%s\": %s", jsonapi.MustEscape(obj.%s))`,
			gooostrings.ToSnakeCase(field.Name),
			"%s",
			field.Name,
		)
		fields = append(
			fields,
			v,
		)
	}
	str := "lines :=  []string{\n"
	str += strings.Join(fields, ", \n") + ",\n"
	str += "}\n"
	str += "return fmt.Sprintf(\"{\\n%s\\n}\", strings.Join(lines, \", \\n\")), nil"

	return template.Method{
		Receiver:    s.Schema.Name,
		Name:        "JSONAPISerialize",
		Args:        []template.Arg{},
		ReturnTypes: []string{"string", "error"},
		Body:        str,
	}.String()
}
