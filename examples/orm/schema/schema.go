package schema

import "github.com/version-1/gooo/pkg/datasource/schema"

var UserSchema = schema.Schema{
	Name:      "User",
	TableName: "users",
	Fields: []schema.Field{
		{
			Name: "ID",
			Type: schema.UUID,
			Options: schema.FieldOptions{
				PrimaryKey: true,
				Immutable:  true,
			},
		},
		{
			Name:    "Username",
			Type:    schema.String,
			Options: schema.FieldOptions{},
		},
		{
			Name:    "Bio",
			Type:    schema.Ref(schema.String),
			Options: schema.FieldOptions{},
		},
		{
			Name:    "Email",
			Type:    schema.String,
			Options: schema.FieldOptions{},
		},
		{
			Name: "CreatedAt",
			Type: schema.Time,
			Options: schema.FieldOptions{
				Immutable: true,
			},
		},
		{
			Name: "UpdatedAt",
			Type: schema.Time,
			Options: schema.FieldOptions{
				Immutable: true,
			},
		},
	},
}
