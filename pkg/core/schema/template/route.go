package template

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/version-1/gooo/pkg/core/schema/openapi/v3_0_0"
	"github.com/version-1/gooo/pkg/core/schema/openapi/yaml"
	"github.com/version-1/gooo/pkg/toolkit/errors"
)

type Route struct {
	InputType  string
	OutputType string
	Method     string
	Path       string
}

func renderRoutes(routes []Route) (string, error) {
	var b bytes.Buffer
	for _, r := range routes {
		tmpl := template.Must(template.New("route").ParseFS(tmpl, "components/route.go.tmpl"))
		if err := tmpl.ExecuteTemplate(&b, "route.go.tmpl", r); err != nil {
			return "", errors.Wrap(err)
		}
	}
	return b.String(), nil
}

func extractRoutes(r *v3_0_0.RootSchema) []Route {
	routes := []Route{}
	r.Paths.Each(func(path string, pathItem v3_0_0.PathItem) error {
		m := map[string]*v3_0_0.Operation{
			"Get":    pathItem.Get,
			"Post":   pathItem.Post,
			"Patch":  pathItem.Patch,
			"Put":    pathItem.Put,
			"Delete": pathItem.Delete,
		}
		for k, v := range m {
			if v == nil {
				continue
			}

			if k == "Get" || k == "Delete" {
				route := Route{
					InputType:  "request.Void",
					OutputType: withSchemaPackageName(detectOutputType(v, 200, "application/json")),
					Method:     k,
					Path:       path,
				}

				routes = append(routes, route)
			} else {
				statusCode := 200
				if k == "Post" {
					statusCode = 201
				}
				route := Route{
					InputType:  withSchemaPackageName(detectInputType(v, "application/json")),
					OutputType: withSchemaPackageName(detectOutputType(v, statusCode, "application/json")),
					Method:     k,
					Path:       path,
				}
				routes = append(routes, route)
			}
		}

		return nil
	})

	return routes
}

func detectInputType(op *v3_0_0.Operation, contentType string) string {
	schema := op.RequestBody.Content.Get(contentType).Schema
	ref := ""
	if schema.Ref != "" {
		ref = schema.Ref
	}

	if schema.Items.Type == "array" && schema.Items.Ref != "" {
		ref = schema.Items.Ref
	}

	schemaName := strings.Replace(ref, "#/components/schemas/", "", 1)
	return schemaName
}

func detectOutputType(op *v3_0_0.Operation, statusCode int, contentType string) string {
	responses := yaml.OrderedMap[v3_0_0.Response](op.Responses)
	schema := responses.Get(strconv.Itoa(statusCode)).Content.Get(contentType).Schema
	ref := ""
	if schema.Ref != "" {
		ref = schema.Ref
	}

	if schema.Type == "array" && schema.Items.Ref != "" {
		ref = schema.Items.Ref
	}

	schemaName := strings.Replace(ref, "#/components/schemas/", "", 1)
	return schemaName
}
