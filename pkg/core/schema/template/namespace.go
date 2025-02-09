package template

import (
	"fmt"
	"strings"
)

func withSchemaPackageName(schemaName string) string {
	return fmt.Sprintf("schema.%s", schemaName)
}

func schemaTypeName(schemaName string) string {
	segments := strings.Split(schemaName, "/")
	return segments[len(segments)-1]
}
