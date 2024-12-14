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

func (obj Profile) Columns() []string {
	return []string{"id", "user_id", "bio", "created_at", "updated_at"}
}

func (obj *Profile) Scan(rows scanner) error {
	if err := rows.Scan(&obj.ID, &obj.UserID, &obj.Bio, &obj.CreatedAt, &obj.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (obj *Profile) Destroy(ctx context.Context, qr queryer) error {
	zero, err := util.IsZero(obj.ID)
	if err != nil {
		return goooerrors.Wrap(err)
	}

	if zero {
		return goooerrors.Wrap(ErrPrimaryKeyMissing)
	}

	query := "DELETE FROM profiles WHERE id = $1"
	if _, err := qr.ExecContext(ctx, query, obj.ID); err != nil {
		return goooerrors.Wrap(err)
	}

	return nil
}

func (obj *Profile) Find(ctx context.Context, qr queryer) error {
	zero, err := util.IsZero(obj.ID)
	if err != nil {
		return goooerrors.Wrap(err)
	}

	if zero {
		return goooerrors.Wrap(ErrPrimaryKeyMissing)
	}

	query := "SELECT id, user_id, bio, created_at, updated_at FROM profiles WHERE id = $1"
	row := qr.QueryRowContext(ctx, query, obj.ID)

	if err := obj.Scan(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goooerrors.Wrap(ErrNotFound)
		}

		return goooerrors.Wrap(err)
	}

	return nil
}

func (obj *Profile) Save(ctx context.Context, qr queryer) error {
	if err := obj.validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO profiles (user_id, bio) VALUES ($1, $2, $3)
		ON CONFLICT(id) DO UPDATE SET user_id = $1, bio = $2, updated_at = NOW()
		RETURNING id, user_id, bio, created_at, updated_at
  `

	row := qr.QueryRowContext(ctx, query, obj.UserID, obj.Bio)
	if err := obj.Scan(row); err != nil {
		return err
	}

	return nil
}

func (obj *Profile) Assign(v Profile) {
	obj.ID = v.ID
	obj.UserID = v.UserID
	obj.Bio = v.Bio
	obj.CreatedAt = v.CreatedAt
	obj.UpdatedAt = v.UpdatedAt
}

func (obj Profile) validate() ormerrors.ValidationError {
	return nil
}

func (obj Profile) JSONAPISerialize() (string, error) {
	lines := []string{
		fmt.Sprintf("\"user_id\": %s", jsonapi.MustEscape(obj.UserID)),
		fmt.Sprintf("\"bio\": %s", jsonapi.MustEscape(obj.Bio)),
		fmt.Sprintf("\"created_at\": %s", jsonapi.MustEscape(obj.CreatedAt)),
		fmt.Sprintf("\"updated_at\": %s", jsonapi.MustEscape(obj.UpdatedAt)),
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(lines, ", \n")), nil
}

func (obj Profile) ToJSONAPIResource() (jsonapi.Resource, jsonapi.Resources) {
	includes := &jsonapi.Resources{ShouldSort: true}
	r := &jsonapi.Resource{
		ID:            jsonapi.Stringify(obj.ID),
		Type:          "profile",
		Attributes:    obj,
		Relationships: jsonapi.Relationships{},
	}

	return *r, *includes
}
