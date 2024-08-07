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

type Post struct {
	schema.Schema
	// db related fields
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	Body      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time

	// non-db related fields
	User User
}

func (obj Post) Columns() []string {
	return []string{"id", "user_id", "title", "body", "status", "created_at", "updated_at"}
}

func (obj *Post) Scan(rows scanner) error {
	if err := rows.Scan(&obj.ID, &obj.UserID, &obj.Title, &obj.Body, &obj.Status, &obj.CreatedAt, &obj.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (obj *Post) Destroy(ctx context.Context, qr queryer) error {
	if obj.ID == uuid.Nil {
		return ErrPrimaryKeyMissing
	}

	query := "DELETE FROM posts WHERE id = $1"
	if _, err := qr.ExecContext(ctx, query, obj.ID); err != nil {
		return err
	}

	return nil
}

func (obj *Post) Find(ctx context.Context, qr queryer) error {
	if obj.ID == uuid.Nil {
		return ErrPrimaryKeyMissing
	}

	query := "SELECT id, user_id, title, body, status, created_at, updated_at FROM posts WHERE id = $1"
	row := qr.QueryRowContext(ctx, query, obj.ID)

	if err := obj.Scan(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (obj *Post) Save(ctx context.Context, qr queryer) error {
	if err := obj.validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO posts (user_id, title, body, status) VALUES ($1, $2, $3)
		ON CONFLICT(id) DO UPDATE SET user_id = $1, title = $2, body = $3, status = $4, updated_at = NOW()
		RETURNING id, user_id, title, body, status, created_at, updated_at
  `

	row := qr.QueryRowContext(ctx, query, obj.UserID, obj.Title, obj.Body, obj.Status)
	if err := obj.Scan(row); err != nil {
		return err
	}

	return nil
}

func (obj *Post) Assign(v Post) {
	obj.ID = v.ID
	obj.UserID = v.UserID
	obj.Title = v.Title
	obj.Body = v.Body
	obj.Status = v.Status
	obj.CreatedAt = v.CreatedAt
	obj.UpdatedAt = v.UpdatedAt
	obj.User = v.User

}

func (obj Post) validate() goooerrors.ValidationError {
	validator := obj.Schema.Fields[1].Options.Validators[0]
	if err := validator.Validate("UserID")(obj.UserID); err != nil {
		return err
	}

	validator = obj.Schema.Fields[2].Options.Validators[0]
	if err := validator.Validate("Title")(obj.Title); err != nil {
		return err
	}

	validator = obj.Schema.Fields[3].Options.Validators[0]
	if err := validator.Validate("Body")(obj.Body); err != nil {
		return err
	}

	return nil
}

func (obj Post) JSONAPISerialize() (string, error) {
	lines := []string{
		fmt.Sprintf("\"id\": %s", jsonapi.Stringify(obj.ID)),
		fmt.Sprintf("\"user_id\": %s", jsonapi.Stringify(obj.UserID)),
		fmt.Sprintf("\"title\": %s", jsonapi.Stringify(obj.Title)),
		fmt.Sprintf("\"body\": %s", jsonapi.Stringify(obj.Body)),
		fmt.Sprintf("\"status\": %s", jsonapi.Stringify(obj.Status)),
		fmt.Sprintf("\"created_at\": %s", jsonapi.Stringify(obj.CreatedAt)),
		fmt.Sprintf("\"updated_at\": %s", jsonapi.Stringify(obj.UpdatedAt)),
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(lines, ", \n")), nil
}

func (obj Post) ToJSONAPIResource() (jsonapi.Resource, jsonapi.Resources) {
	includes := &jsonapi.Resources{}
	r := jsonapi.Resource{
		ID:            jsonapi.Stringify(obj.ID),
		Type:          "post",
		Attributes:    obj,
		Relationships: jsonapi.Relationships{},
	}

	ele := obj.User
	if ele.ID == (User{}).ID {
		return r, *includes
	}
	relationship := jsonapi.Relationship{
		Data: jsonapi.ResourceIdentifier{
			ID:   jsonapi.Stringify(ele.ID),
			Type: "user",
		},
	}

	resource, childIncludes := ele.ToJSONAPIResource()
	includes.Append(resource)
	includes.Append(childIncludes.Data...)

	r.Relationships["user"] = relationship
	return r, *includes
}
