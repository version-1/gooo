package route

import (
	"strconv"
	gostrings "strings"

	"github.com/version-1/gooo/pkg/core/schema/openapi/v3_0_0"
	"github.com/version-1/gooo/pkg/core/schema/openapi/yaml"
	"github.com/version-1/gooo/pkg/toolkit/strings"
)

func ResolveRouteName(method, path string, operation *v3_0_0.Operation) string {
	if operation.OperationId != "" {
		return strings.ToCamelCase(operation.OperationId)
	}

	name := strings.ToPascalCase(method)
	for _, p := range gostrings.Split(path, "/") {
		if gostrings.HasPrefix(p, "{") && gostrings.HasSuffix(p, "}") {
			name += strings.ToPascalCase(p[1 : len(p)-1])
		} else {
			name += strings.ToPascalCase(p)
		}
	}

	return name
}

func DetectInputType(op *v3_0_0.Operation, contentType string) string {
	schema := op.RequestBody.Content.Get(contentType).Schema
	ref := ""
	if schema.Ref != "" {
		ref = schema.Ref
	}

	if schema.Items.Type == "array" && schema.Items.Ref != "" {
		ref = schema.Items.Ref
	}

	schemaName := gostrings.Replace(ref, "#/components/schemas/", "", 1)
	return schemaName
}

func DetectOutputType(op *v3_0_0.Operation, statusCode int, contentType string) string {
	responses := yaml.OrderedMap[v3_0_0.Response](op.Responses)
	schema := responses.Get(strconv.Itoa(statusCode)).Content.Get(contentType).Schema
	ref := ""
	if schema.Ref != "" {
		ref = schema.Ref
	}

	if schema.Type == "array" && schema.Items.Ref != "" {
		ref = schema.Items.Ref
	}

	schemaName := gostrings.Replace(ref, "#/components/schemas/", "", 1)
	return schemaName
}
