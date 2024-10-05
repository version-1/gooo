package fixtures

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	ormerrors "github.com/version-1/gooo/pkg/datasource/orm/errors"
	goooerrors "github.com/version-1/gooo/pkg/errors"
	"github.com/version-1/gooo/pkg/presenter/jsonapi"
	"github.com/version-1/gooo/pkg/util"
)

func (obj Post) Columns() []string {
	return []string{"id", "user_id", "title", "body", "created_at", "updated_at"}
}

func (obj *Post) Scan(rows scanner) error {
	if err := rows.Scan(&obj.ID, &obj.UserID, &obj.Title, &obj.Body, &obj.CreatedAt, &obj.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (obj *Post) Destroy(ctx context.Context, qr queryer) error {
	zero, err := util.IsZero(obj.ID)
	if err != nil {
		return goooerrors.Wrap(err)
	}

	if zero {
		return goooerrors.Wrap(ErrPrimaryKeyMissing)
	}

	query := "DELETE FROM posts WHERE id = $1"
	if _, err := qr.ExecContext(ctx, query, obj.ID); err != nil {
		return goooerrors.Wrap(err)
	}

	return nil
}

func (obj *Post) Find(ctx context.Context, qr queryer) error {
	zero, err := util.IsZero(obj.ID)
	if err != nil {
		return goooerrors.Wrap(err)
	}

	if zero {
		return goooerrors.Wrap(ErrPrimaryKeyMissing)
	}

	query := "SELECT id, user_id, title, body, created_at, updated_at FROM posts WHERE id = $1"
	row := qr.QueryRowContext(ctx, query, obj.ID)

	if err := obj.Scan(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goooerrors.Wrap(ErrNotFound)
		}

		return goooerrors.Wrap(err)
	}

	return nil
}

func (obj *Post) Save(ctx context.Context, qr queryer) error {
	if err := obj.validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO posts (user_id, title, body, user, likes) VALUES ($1, $2, $3)
		ON CONFLICT(id) DO UPDATE SET user_id = $1, title = $2, body = $3, user = $4, likes = $5, updated_at = NOW()
		RETURNING id, user_id, title, body, created_at, updated_at
  `

	row := qr.QueryRowContext(ctx, query, obj.UserID, obj.Title, obj.Body, obj.User, obj.Likes)
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
	obj.CreatedAt = v.CreatedAt
	obj.UpdatedAt = v.UpdatedAt
	obj.User = v.User
	obj.Likes = v.Likes
}

func (obj Post) validate() ormerrors.ValidationError {
	return nil
}

func (obj Post) JSONAPISerialize() (string, error) {
	lines := []string{
		fmt.Sprintf("\"user_id\": %s", jsonapi.MustEscape(obj.UserID)),
		fmt.Sprintf("\"title\": %s", jsonapi.MustEscape(obj.Title)),
		fmt.Sprintf("\"body\": %s", jsonapi.MustEscape(obj.Body)),
		fmt.Sprintf("\"created_at\": %s", jsonapi.MustEscape(obj.CreatedAt)),
		fmt.Sprintf("\"updated_at\": %s", jsonapi.MustEscape(obj.UpdatedAt)),
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(lines, ", \n")), nil
}

func (obj Post) ToJSONAPIResource() (jsonapi.Resource, jsonapi.Resources) {
	includes := &jsonapi.Resources{ShouldSort: true}
	r := &jsonapi.Resource{
		ID:            jsonapi.Stringify(obj.ID),
		Type:          "post",
		Attributes:    obj,
		Relationships: jsonapi.Relationships{},
	}

	ele := obj.User
	if ele.ID != (User{}).ID {
		jsonapi.HasOne(r, includes, ele, ele.ID, "user")
	}

	elements := []jsonapi.Resourcer{}
	for _, ele := range obj.Likes {
		elements = append(elements, jsonapi.Resourcer(ele))
	}
	jsonapi.HasMany(r, includes, elements, "like", func(ri *jsonapi.ResourceIdentifier, i int) {
		id := obj.Likes[i].ID
		ri.ID = jsonapi.Stringify(id)
	})

	return *r, *includes
}
