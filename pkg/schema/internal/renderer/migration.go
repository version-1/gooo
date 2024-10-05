package renderer

import (
	"fmt"
	"strings"
)

type migrationSchema interface {
	SchemaNames() []string
}

type InitialMigraiton struct {
	Path   string
	schema migrationSchema
}

func NewInitialMigration(path string, s migrationSchema) *InitialMigraiton {
	return &InitialMigraiton{
		Path:   path,
		schema: s,
	}
}

func (i InitialMigraiton) Filename() string {
	return fmt.Sprintf("%s-initial.yaml", strings.Repeat("0", 14))
}

func (i InitialMigraiton) Render() (string, error) {
	for _, name := range i.schema.SchemaNames() {
		fmt.Println(name)
	}
	return "", nil
}
