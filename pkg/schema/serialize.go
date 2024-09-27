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
		r := jsonapi.Resource{
		  ID:   jsonapi.Stringify(obj.%s),
			Type: "%s",
			Attributes: obj,
			Relationships: jsonapi.Relationships{},
		}

	`, primaryKey, gooostrings.ToSnakeCase(s.Schema.Name))

	for _, field := range s.Schema.AssociationFields() {
		t := fmt.Stringer(field.Type)
		_, ok := t.(slice)
		if v, ok := t.(Elementer); ok {
			t = v.Element()
		}
		typeName := gooostrings.ToSnakeCase(t.String())
		association := field.Options.Association
		primaryKey := association.Schema.PrimaryKey()
		if ok {
			str += fmt.Sprintf(`
			relationships := jsonapi.RelationshipHasMany{}
			for _, ele := range obj.%s {
				relationships.Data = append(
					relationships.Data,
					jsonapi.ResourceIdentifier{
						ID:   jsonapi.Stringify(ele.%s),
						Type: "%s",
					},
				)

				resource, childIncludes := ele.ToJSONAPIResource()
				includes.Append(resource)
				includes.Append(childIncludes.Data...)
			}

			if len(relationships.Data) > 0 {
				r.Relationships["%s"] = relationships
			}
		`, field.Name, primaryKey, typeName, gooostrings.ToSnakeCase(field.Name))
		} else {
			str += fmt.Sprintf(`
			ele := obj.%s
			if ele.%s == (%s{}).%s {
				return r, *includes
			}
			relationship := jsonapi.Relationship{
				Data: jsonapi.ResourceIdentifier{
					ID:   jsonapi.Stringify(ele.%s),
					Type: "%s",
				},
			}

			resource, childIncludes := ele.ToJSONAPIResource()
			includes.Append(resource)
			includes.Append(childIncludes.Data...)

			r.Relationships["%s"] = relationship
		`, field.Name, primaryKey, t.String(), primaryKey, primaryKey, typeName, typeName)
		}
	}

	str += "return r, *includes"

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
	for _, field := range s.Schema.ColumnFields() {
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
