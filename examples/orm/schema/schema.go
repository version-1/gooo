package schema

import (
	"strings"

	"github.com/version-1/gooo/pkg/datasource/orm/errors"
	"github.com/version-1/gooo/pkg/datasource/orm/validator"
	"github.com/version-1/gooo/pkg/datasource/schema"
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
