package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/version-1/gooo/examples/orm/schema"
	"github.com/version-1/gooo/pkg/datasource/orm"
)

func main() {
	args := os.Args[1:]

	dirpath := args[0]
	schema := orm.SchemaCollection{
		URL:     "github.com/version-1/gooo",
		Dir:     dirpath,
		Package: filepath.Base(dirpath),
		Schemas: []orm.Schema{
			schema.UserSchema,
		},
	}

	if err := schema.Gen(); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}
