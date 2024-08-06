package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	goooerrors "github.com/version-1/gooo/pkg/datasource/orm/errors"
	"github.com/version-1/gooo/pkg/presenter/jsonapi"
	"github.com/version-1/gooo/pkg/schema"
)

type User struct {
	schema.Schema
	// db related fields
	ID        uuid.UUID
	Username  string
	Bio       *string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time

	// non-db related fields
	Posts []Post
}

func (obj User) Columns() []string {
	return []string{"id", "username", "bio", "email", "created_at", "updated_at"}
}

func (obj *User) Scan(rows scanner) error {
	if err := rows.Scan(&obj.ID, &obj.Username, &obj.Bio, &obj.Email, &obj.CreatedAt, &obj.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (obj *User) Destroy(ctx context.Context, qr queryer) error {
	if obj.ID == uuid.Nil {
		return ErrPrimaryKeyMissing
	}

	query := "DELETE FROM users WHERE id = $1"
	if _, err := qr.ExecContext(ctx, query, obj.ID); err != nil {
		return err
	}

	return nil
}

func (obj *User) Find(ctx context.Context, qr queryer) error {
	if obj.ID == uuid.Nil {
		return ErrPrimaryKeyMissing
	}

	query := "SELECT id, username, bio, email, created_at, updated_at FROM users WHERE id = $1"
	row := qr.QueryRowContext(ctx, query, obj.ID)

	if err := obj.Scan(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (obj *User) Save(ctx context.Context, qr queryer) error {
	if err := obj.validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO users (username, bio, email) VALUES ($1, $2, $3)
		ON CONFLICT(id) DO UPDATE SET username = $1, bio = $2, email = $3, updated_at = NOW()
		RETURNING id, username, bio, email, created_at, updated_at
  `

	row := qr.QueryRowContext(ctx, query, obj.Username, obj.Bio, obj.Email)
	if err := obj.Scan(row); err != nil {
		return err
	}

	return nil
}

func (obj *User) Assign(v User) {
	obj.ID = v.ID
	obj.Username = v.Username
	obj.Bio = v.Bio
	obj.Email = v.Email
	obj.CreatedAt = v.CreatedAt
	obj.UpdatedAt = v.UpdatedAt
	obj.Posts = v.Posts

}

func (obj User) validate() goooerrors.ValidationError {
	validator := obj.Schema.Fields[1].Options.Validators[0]
	if err := validator.Validate("Username")(obj.Username); err != nil {
		return err
	}

	validator = obj.Schema.Fields[1].Options.Validators[1]
	if err := validator.Validate("Username")(obj.Username, obj.Email); err != nil {
		return err
	}

	validator = obj.Schema.Fields[3].Options.Validators[0]
	if err := validator.Validate("Email")(obj.Email); err != nil {
		return err
	}

	validator = obj.Schema.Fields[3].Options.Validators[1]
	if err := validator.Validate("Email")(obj.Email); err != nil {
		return err
	}

	return nil
}

func (obj User) JSONAPISerialize() (string, error) {
	lines := []string{
		fmt.Sprintf("\"id\": %s", jsonapi.Stringify(obj.ID)),
		fmt.Sprintf("\"username\": %s", jsonapi.Stringify(obj.Username)),
		fmt.Sprintf("\"bio\": %s", jsonapi.Stringify(obj.Bio)),
		fmt.Sprintf("\"email\": %s", jsonapi.Stringify(obj.Email)),
		fmt.Sprintf("\"created_at\": %s", jsonapi.Stringify(obj.CreatedAt)),
		fmt.Sprintf("\"updated_at\": %s", jsonapi.Stringify(obj.UpdatedAt)),
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(lines, ", \n")), nil
}

func (obj User) ToJSONAPIResource() (jsonapi.Resource, jsonapi.Resources) {
	includes := &jsonapi.Resources{}
	r := jsonapi.Resource{
		ID:            jsonapi.Stringify(obj.ID),
		Type:          "user",
		Attributes:    obj,
		Relationships: jsonapi.Relationships{},
	}

	relationships := jsonapi.RelationshipHasMany{}
	for _, ele := range obj.Posts {
		relationships.Data = append(
			relationships.Data,
			jsonapi.ResourceIdentifier{
				ID:   jsonapi.Stringify(ele.ID),
				Type: "post",
			},
		)

		resource, childIncludes := ele.ToJSONAPIResource()
		includes.Append(resource)
		includes.Append(childIncludes.Data...)
	}

	if len(relationships.Data) > 0 {
		r.Relationships["posts"] = relationships
	}
	return r, *includes
}
