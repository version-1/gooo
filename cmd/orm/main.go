package main

import (
	"fmt"
	"os"
	"path/filepath"

	exampleschema "github.com/version-1/gooo/examples/orm/schema"
	goooschema "github.com/version-1/gooo/pkg/datasource/schema"
)

func main() {
	args := os.Args[1:]

	dirpath := args[0]
	schema := goooschema.SchemaCollection{
		URL:     "github.com/version-1/gooo",
		Dir:     dirpath,
		Package: filepath.Base(dirpath),
		Schemas: []goooschema.Schema{
			exampleschema.UserSchema,
		},
	}

	if err := schema.Gen(); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}
