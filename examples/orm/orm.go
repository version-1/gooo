package main

import (
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/version-1/gooo/pkg/datasource/orm"
	"github.com/version-1/gooo/pkg/datasource/orm/errors"
	"github.com/version-1/gooo/pkg/datasource/orm/validator"
)

var defaultSchema = &orm.SchemaFactory{
	Primary: orm.Field{
		Name: "ID",
		Options: orm.FieldOptions{
			Immutable: true,
		},
	},
	DefaultFields: []orm.Field{
		{
			Name: "CreatedAt",
			Options: orm.FieldOptions{
				Immutable: true,
			},
		},
		{
			Name: "UpdatedAt",
			Options: orm.FieldOptions{
				Immutable: true,
			},
		},
	},
}

type userSchema struct {
	schema orm.Schema
}

var UserSchema = &userSchema{
	*defaultSchema.NewSchema([]orm.Field{
		{
			Name: "Username",
			Options: orm.FieldOptions{
				Validators: []validator.ValidateFunc{
					validator.Required("Username"),
				},
			},
		},
		{
			Name: "Email",
			Options: orm.FieldOptions{
				Validators: []validator.ValidateFunc{
					validator.Required("Email"),
				},
			},
		},
	}),
}

func (u userSchema) Scan(s orm.Scanner) (orm.Model, error) {
	v := &user{}

	if err := s.Scan(&v.ID, &v.Username, &v.Email, &v.CreatedAt, &v.UpdatedAt); err != nil {
		return v, err
	}

	return v, nil
}

type user struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var _ orm.Model = user{}

func (u user) Scan(s orm.Scanner) (orm.Model, error) {
	return UserSchema.Scan(s)
}

func (u user) Validate() errors.ValidationError {
	for i, f := range UserSchema.schema.Fields {
		for j, validator := range f.Validators {
			if err := validator(u.Values()[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (u user) Fields() []string {
	return UserSchema.schema.FieldKeys()
}

func (u user) MutableFields() []string {
	return UserSchema.schema.MutableFieldKeys()
}

func (u user) Values() []any {
	return []any{u.Username, u.Email}
}

func (u user) TableName() string {
	return "users"
}

func (u user) Identifier() string {
	return u.ID.String()
}

func (u user) NewItem() orm.Model {
	return &user{}
}

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	ormFactory := orm.NewOrmFactory(db, log.New(os.Stdout, "", 0), orm.Options{QueryLog: true})

	o := ormFactory.New(UserSchema)
}
