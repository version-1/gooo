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

func (obj User) Columns() []string {
	return []string{"id", "username", "email", "refresh_token", "timezone", "time_diff", "created_at", "updated_at"}
}

func (obj *User) Scan(rows scanner) error {
	if err := rows.Scan(&obj.ID, &obj.Username, &obj.Email, &obj.RefreshToken, &obj.Timezone, &obj.TimeDiff, &obj.CreatedAt, &obj.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (obj *User) Destroy(ctx context.Context, qr queryer) error {
	zero, err := util.IsZero(obj.ID)
	if err != nil {
		return goooerrors.Wrap(err)
	}

	if zero {
		return goooerrors.Wrap(ErrPrimaryKeyMissing)
	}

	query := "DELETE FROM users WHERE id = $1"
	if _, err := qr.ExecContext(ctx, query, obj.ID); err != nil {
		return goooerrors.Wrap(err)
	}

	return nil
}

func (obj *User) Find(ctx context.Context, qr queryer) error {
	zero, err := util.IsZero(obj.ID)
	if err != nil {
		return goooerrors.Wrap(err)
	}

	if zero {
		return goooerrors.Wrap(ErrPrimaryKeyMissing)
	}

	query := "SELECT id, username, email, refresh_token, timezone, time_diff, created_at, updated_at FROM users WHERE id = $1"
	row := qr.QueryRowContext(ctx, query, obj.ID)

	if err := obj.Scan(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goooerrors.Wrap(ErrNotFound)
		}

		return goooerrors.Wrap(err)
	}

	return nil
}

func (obj *User) Save(ctx context.Context, qr queryer) error {
	if err := obj.validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO users (username, email, refresh_token, timezone, time_diff, profile, posts) VALUES ($1, $2, $3)
		ON CONFLICT(id) DO UPDATE SET username = $1, email = $2, refresh_token = $3, timezone = $4, time_diff = $5, profile = $6, posts = $7, updated_at = NOW()
		RETURNING id, username, email, refresh_token, timezone, time_diff, created_at, updated_at
  `

	row := qr.QueryRowContext(ctx, query, obj.Username, obj.Email, obj.RefreshToken, obj.Timezone, obj.TimeDiff, obj.Profile, obj.Posts)
	if err := obj.Scan(row); err != nil {
		return err
	}

	return nil
}

func (obj *User) Assign(v User) {
	obj.ID = v.ID
	obj.Username = v.Username
	obj.Email = v.Email
	obj.RefreshToken = v.RefreshToken
	obj.Timezone = v.Timezone
	obj.TimeDiff = v.TimeDiff
	obj.CreatedAt = v.CreatedAt
	obj.UpdatedAt = v.UpdatedAt
	obj.Profile = v.Profile
	obj.Posts = v.Posts
}

func (obj User) validate() ormerrors.ValidationError {
	return nil
}

func (obj User) JSONAPISerialize() (string, error) {
	lines := []string{
		fmt.Sprintf("\"username\": %s", jsonapi.MustEscape(obj.Username)),
		fmt.Sprintf("\"email\": %s", jsonapi.MustEscape(obj.Email)),
		fmt.Sprintf("\"refresh_token\": %s", jsonapi.MustEscape(obj.RefreshToken)),
		fmt.Sprintf("\"timezone\": %s", jsonapi.MustEscape(obj.Timezone)),
		fmt.Sprintf("\"time_diff\": %s", jsonapi.MustEscape(obj.TimeDiff)),
		fmt.Sprintf("\"created_at\": %s", jsonapi.MustEscape(obj.CreatedAt)),
		fmt.Sprintf("\"updated_at\": %s", jsonapi.MustEscape(obj.UpdatedAt)),
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(lines, ", \n")), nil
}

func (obj User) ToJSONAPIResource() (jsonapi.Resource, jsonapi.Resources) {
	includes := &jsonapi.Resources{ShouldSort: true}
	r := &jsonapi.Resource{
		ID:            jsonapi.Stringify(obj.ID),
		Type:          "user",
		Attributes:    obj,
		Relationships: jsonapi.Relationships{},
	}

	ele := obj.Profile
	if ele != nil {
		jsonapi.HasOne(r, includes, ele, ele.ID, "profile")
	}

	elements := []jsonapi.Resourcer{}
	for _, ele := range obj.Posts {
		elements = append(elements, jsonapi.Resourcer(ele))
	}
	jsonapi.HasMany(r, includes, elements, "post", func(ri *jsonapi.ResourceIdentifier, i int) {
		id := obj.Posts[i].ID
		ri.ID = jsonapi.Stringify(id)
	})

	return *r, *includes
}
