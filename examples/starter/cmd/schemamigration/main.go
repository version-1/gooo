package main

import (
	"fmt"
	"strings"

	"github.com/version-1/gooo/pkg/schema"
)

func main() {
	p := schema.NewParser()
	list, err := p.Parse("./pkg/schema/internal/schema/schema.go")
	if err != nil {
		panic(err)
	}

	s := schema.SchemaCollection{
		URL:     "github.com/version-1/gooo",
		Dir:     "internal/schema",
		Package: "schema",
		Schemas: list,
	}

	m := schema.NewMigration(s, schema.MigrationConfig{})
	os, err := m.OriginSchema()
	if err != nil {
		panic(err)
	}

	filename := fmt.Sprintf("%s_initital.yaml", strings.Repeat("0", 14))
	path := fmt.Sprintf("examples/starter/db/v2/migrations/%s", filename)
	fmt.Printf("Writing to %s\n", path)
	if err := os.Write(path); err != nil {
		fmt.Printf("Error: %+v\n", err)
		panic(err)
	}
}
