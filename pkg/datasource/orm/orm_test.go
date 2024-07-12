package orm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	ormerrors "github.com/version-1/gooo/pkg/datasource/orm/errors"
	"github.com/version-1/gooo/pkg/datasource/orm/validator"
)

var defaultSchema = &SchemaFactory{
	Primary: Field{
		Name: "ID",
		Options: FieldOptions{
			Immutable: true,
		},
	},
	DefaultFields: []Field{
		{
			Name: "CreatedAt",
			Options: FieldOptions{
				Immutable: true,
			},
		},
		{
			Name: "UpdatedAt",
			Options: FieldOptions{
				Immutable: true,
			},
		},
	},
}

type userSchema struct {
	schema Schema
}

var UserSchema = &userSchema{
	*defaultSchema.NewSchema([]Field{
		{
			Name: "Username",
			Options: FieldOptions{
				Validators: []validator.ValidateFunc{
					validator.Required("Username"),
				},
			},
		},
		{
			Name: "Email",
			Options: FieldOptions{
				Validators: []validator.ValidateFunc{
					validator.Required("Email"),
				},
			},
		},
	}),
}

func (u userSchema) Scan(s Scanner) (Model, error) {
	v := &user{}

	if err := s.Scan(&v.ID, &v.Username, &v.Email, &v.CreatedAt, &v.UpdatedAt); err != nil {
		return v, err
	}

	return v, nil
}

type testLogger struct {
	messages [][]string
}

var _ Logger = &testLogger{}

func (l *testLogger) Warnf(format string, args ...interface{}) {
	l.messages = append(l.messages, []string{"warn", fmt.Sprintf(format, args...)})
}

func (l *testLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	l.messages = append(l.messages, []string{"info", fmt.Sprintf(format, args...)})
}

func (l *testLogger) Debugf(format string, args ...interface{}) {
	l.messages = append(l.messages, []string{"debug", fmt.Sprintf(format, args...)})
}

func (l *testLogger) Errorf(format string, args ...interface{}) {
	l.messages = append(l.messages, []string{"error", fmt.Sprintf(format, args...)})
}

type user struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var _ Model = user{}

func (u user) Scan(s Scanner) (Model, error) {
	return UserSchema.Scan(s)
}

func (u user) Validate() ormerrors.ValidationError {
	return UserSchema.schema.Validate(u)
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

func (u user) NewItem() Model {
	return &user{}
}

func TestOperators(t *testing.T) {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	ormFactory := NewOrmFactory(db, &testLogger{}, Options{QueryLog: true})

	o := ormFactory.New(UserSchema.schema)

	u := user{
		Username: "gooo",
		Email:    "gooo@example.com",
	}

	// create
	if err := Create[*user](context.Background(), o, &u); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("User: %#v\n", u)

	if u.ID == uuid.Nil {
		t.Fatal("ID is not set")
	}

	if u.Username != "gooo" {
		t.Fatal("Name is not set")
	}

	if u.Email != "gooo@example.com" {
		t.Fatal("Email is not set")
	}

	if u.CreatedAt.IsZero() {
		t.Fatal("CreatedAt is not set")
	}

	if u.UpdatedAt.IsZero() {
		t.Fatal("UpdatedAt is not set")
	}

	// find
	uu := user{ID: u.ID}
	if err := Find[*user](context.Background(), o, &uu); err != nil {
		t.Fatal(err)
	}

	if u.ID == uuid.Nil {
		t.Fatal("ID is not set")
	}

	if u.Username != "gooo" {
		t.Fatal("Username is not set")
	}

	if u.Email != "gooo@example.com" {
		t.Fatal("Email is not set")
	}

	if u.CreatedAt.IsZero() {
		t.Fatal("CreatedAt is not set")
	}

	if u.UpdatedAt.IsZero() {
		t.Fatal("UpdatedAt is not set")
	}

	// update
	u.Username = "editedgooo"
	u.Email = "editedGooo"
	prevUpdatedAt := u.UpdatedAt
	if err := Update[*user](context.Background(), o, &u); err != nil {
		t.Fatal(err)
	}

	if u.ID == uuid.Nil {
		t.Fatal("ID is not set")
	}

	if u.Username != "editedgooo" {
		t.Fatal("Username is not updated")
	}

	if u.Email != "editedgooo" {
		t.Fatal("Email is not updated")
	}

	if u.UpdatedAt.Equal(prevUpdatedAt) {
		t.Fatal("UpdatedAt is not updated")
	}

	// delete
	if err := Delete[*user](context.Background(), o, &u); err != nil {
		t.Fatal(err)
	}

	if err := Find[*user](context.Background(), o, &uu); err != nil {
		if errors.Is(err, ormerrors.ErrNotFound) {
			return
		}

		t.Fatal(err)
	}
}
