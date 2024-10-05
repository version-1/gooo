package renderer

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
	`, primaryKey, gooostrings.ToSnakeCase(s.Schema.GetName()))
	str += "\n"

	for _, ident := range s.Schema.AssociationFieldIdents() {
		if ident.Slice {
			str += fmt.Sprintf(
				`elements := []jsonapi.Resourcer{}
				for _, ele := range obj.%s {
					elements = append(elements, jsonapi.Resourcer(ele))
				}
				jsonapi.HasMany(r, includes, elements, "%s", func(ri *jsonapi.ResourceIdentifier, i int) {
						id := obj.%s[i].%s
						ri.ID = jsonapi.Stringify(id)
				})`,
				ident.FieldName,
				ident.TypeName,
				ident.FieldName,
				ident.PrimaryKey,
			)
			str += "\n"
		} else {
			if ident.Ref {
				str += fmt.Sprintf(
					`ele := obj.%s
					if ele != nil {
						jsonapi.HasOne(r, includes, ele, ele.%s, "%s")
					}`,
					ident.FieldName,
					ident.PrimaryKey,
					ident.TypeName,
				)
			} else {
				str += fmt.Sprintf(
					`ele := obj.%s
					if ele.%s != (%s{}).%s {
						jsonapi.HasOne(r, includes, ele, ele.%s, "%s")
					}`,
					ident.FieldName,
					ident.PrimaryKey,
					ident.TypeElementExpr,
					ident.PrimaryKey,
					ident.PrimaryKey,
					ident.TypeName,
				)
			}
			str += "\n"
		}
		str += "\n"
	}

	str += "\n"
	str += "return *r, *includes"

	return template.Method{
		Receiver:    s.Schema.GetName(),
		Name:        "ToJSONAPIResource",
		Args:        []template.Arg{},
		ReturnTypes: []string{"jsonapi.Resource", "jsonapi.Resources"},
		Body:        str,
	}.String()
}

func (s SchemaTemplate) defineJSONAPISerialize() string {
	fields := []string{}
	for _, n := range s.Schema.AttributeFieldNames() {
		v := fmt.Sprintf(
			`fmt.Sprintf("\"%s\": %s", jsonapi.MustEscape(obj.%s))`,
			gooostrings.ToSnakeCase(n),
			"%s",
			n,
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
		Receiver:    s.Schema.GetName(),
		Name:        "JSONAPISerialize",
		Args:        []template.Arg{},
		ReturnTypes: []string{"string", "error"},
		Body:        str,
	}.String()
}
