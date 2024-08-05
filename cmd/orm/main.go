package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/version-1/gooo/pkg/datasource/orm"
)

func main() {
	args := os.Args[1:]

	dirpath := args[0]
	schema := orm.SchemaCollection{
		Dir:     dirpath,
		Package: filepath.Base(dirpath),
		Schemas: []orm.Schema{
			user,
		},
	}

	if err := schema.Gen(); err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}

var user = orm.Schema{
	Name:      "User",
	TableName: "users",
	Fields: []orm.Field{
		{
			Name: "ID",
			Type: orm.UUID,
			Options: orm.FieldOptions{
				Immutable:  true,
				PrimaryKey: true,
			},
		},
		{
			Name:    "Username",
			Type:    orm.String,
			Options: orm.FieldOptions{},
		},
		{
			Name: "Bio",
			Type: orm.Ref(orm.String),
		},
		{
			Name:    "Email",
			Type:    orm.String,
			Options: orm.FieldOptions{},
		},
		{
			Name: "CreatedAt",
			Type: orm.Time,
			Options: orm.FieldOptions{
				Immutable: true,
			},
		},
		{
			Name: "UpdatedAt",
			Type: orm.Time,
			Options: orm.FieldOptions{
				Immutable: true,
			},
		},
	},
}
