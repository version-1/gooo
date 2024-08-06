package schema

import "github.com/version-1/gooo/pkg/datasource/orm"

var UserSchema = orm.Schema{
	Name:      "User",
	TableName: "users",
	Fields: []orm.Field{
		{
			Name: "ID",
			Type: orm.UUID,
			Options: orm.FieldOptions{
				PrimaryKey: true,
				Immutable:  true,
			},
		},
		{
			Name:    "Username",
			Type:    orm.String,
			Options: orm.FieldOptions{},
		},
		{
			Name:    "Bio",
			Type:    orm.Ref(orm.String),
			Options: orm.FieldOptions{},
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
