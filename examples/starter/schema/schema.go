package schema

import (
	"path/filepath"

	"github.com/version-1/gooo/pkg/schema"
)

var UserSchema = schema.Schema{
	Name:      "User",
	TableName: "users",
	Fields: []schema.Field{
		{
			Name: "ID",
			Type: schema.UUID,
		},
		{
			Name: "Username",
			Type: schema.String,
		},
		{
			Name: "Bio",
			Type: schema.Ref(schema.String),
		},
		{
			Name: "Email",
			Type: schema.String,
		},
		{
			Name: "CreatedAt",
			Type: schema.Time,
		},
		{
			Name: "UpdatedAt",
			Type: schema.Time,
		},
	},
}

var PostSchema = schema.Schema{
	Name:      "Post",
	TableName: "posts",
	Fields: []schema.Field{
		{
			Name: "ID",
			Type: schema.UUID,
		},
		{
			Name: "UserID",
			Type: schema.UUID,
		},
		{
			Name: "Title",
			Type: schema.String,
		},
		{
			Name: "Body",
			Type: schema.String,
		},
		{
			Name: "Status",
			Type: schema.String,
		},
		{
			Name: "CreatedAt",
			Type: schema.Time,
		},
		{
			Name: "UpdatedAt",
			Type: schema.Time,
		},
	},
}

func Run(dir string) error {
	UserSchema.AddFields(schema.Field{
		Name: "Posts",
		Type: schema.Slice(PostSchema.Type()),
	})

	PostSchema.AddFields(schema.Field{
		Name: "User",
		Type: UserSchema.Type(),
	})

	schemas := schema.SchemaCollection{
		URL:     "github.com/version-1/gooo",
		Package: filepath.Base(dir),
		Dir:     dir,
		Schemas: []schema.Schema{
			UserSchema,
			PostSchema,
		},
	}

	if err := schemas.Gen(); err != nil {
		return err
	}

	return nil
}
