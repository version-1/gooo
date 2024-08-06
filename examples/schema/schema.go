package schema

import (
	"path/filepath"
	"strings"

	"github.com/version-1/gooo/pkg/datasource/orm/errors"
	"github.com/version-1/gooo/pkg/datasource/orm/validator"
	"github.com/version-1/gooo/pkg/schema"
)

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
			Name: "Username",
			Type: schema.String,
			Options: schema.FieldOptions{
				Validators: []schema.Validator{
					{
						Validate: validator.Required,
					},
					{
						Fields: []string{"Email"},
						Validate: func(key string) validator.ValidatorFunc {
							return func(v ...any) errors.ValidationError {
								username := v[0].(string)
								email := strings.Split(v[1].(string), "@")[0]
								if strings.Contains(username, email) {
									return errors.NewValidationError(key, "Username should not contain email")
								}

								return nil
							}
						},
					},
				},
			},
		},
		{
			Name:    "Bio",
			Type:    schema.Ref(schema.String),
			Options: schema.FieldOptions{},
		},
		{
			Name: "Email",
			Type: schema.String,
			Options: schema.FieldOptions{
				Validators: []schema.Validator{
					{
						Validate: validator.Required,
					},
					{
						Validate: validator.Email,
					},
				},
			},
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

var PostSchema = schema.Schema{
	Name:      "Post",
	TableName: "posts",
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
			Name: "UserID",
			Type: schema.UUID,
			Options: schema.FieldOptions{
				Validators: []schema.Validator{
					{
						Validate: validator.Required,
					},
				},
			},
		},
		{
			Name: "Title",
			Type: schema.String,
			Options: schema.FieldOptions{
				Validators: []schema.Validator{
					{
						Validate: validator.Required,
					},
				},
			},
		},
		{
			Name: "Body",
			Type: schema.String,
			Options: schema.FieldOptions{
				Validators: []schema.Validator{
					{
						Validate: validator.Required,
					},
				},
			},
		},
		{
			Name:    "Status",
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

func Run(dir string) error {
	UserSchema.AddFields(schema.Field{
		Name: "Posts",
		Type: schema.Slice(PostSchema.Type()),
		Options: schema.FieldOptions{
			Ignore: true,
		},
	})

	PostSchema.AddFields(schema.Field{
		Name: "User",
		Type: UserSchema.Type(),
		Options: schema.FieldOptions{
			Ignore: true,
		},
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
