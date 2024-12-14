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

func (obj Like) Columns() []string {
	return []string{"id", "likeable_id", "likeable_type", "created_at", "updated_at"}
}

func (obj *Like) Scan(rows scanner) error {
	if err := rows.Scan(&obj.ID, &obj.LikeableID, &obj.LikeableType, &obj.CreatedAt, &obj.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (obj *Like) Destroy(ctx context.Context, qr queryer) error {
	zero, err := util.IsZero(obj.ID)
	if err != nil {
		return goooerrors.Wrap(err)
	}

	if zero {
		return goooerrors.Wrap(ErrPrimaryKeyMissing)
	}

	query := "DELETE FROM likes WHERE id = $1"
	if _, err := qr.ExecContext(ctx, query, obj.ID); err != nil {
		return goooerrors.Wrap(err)
	}

	return nil
}

func (obj *Like) Find(ctx context.Context, qr queryer) error {
	zero, err := util.IsZero(obj.ID)
	if err != nil {
		return goooerrors.Wrap(err)
	}

	if zero {
		return goooerrors.Wrap(ErrPrimaryKeyMissing)
	}

	query := "SELECT id, likeable_id, likeable_type, created_at, updated_at FROM likes WHERE id = $1"
	row := qr.QueryRowContext(ctx, query, obj.ID)

	if err := obj.Scan(row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goooerrors.Wrap(ErrNotFound)
		}

		return goooerrors.Wrap(err)
	}

	return nil
}

func (obj *Like) Save(ctx context.Context, qr queryer) error {
	if err := obj.validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO likes (likeable_id, likeable_type) VALUES ($1, $2, $3)
		ON CONFLICT(id) DO UPDATE SET likeable_id = $1, likeable_type = $2, updated_at = NOW()
		RETURNING id, likeable_id, likeable_type, created_at, updated_at
  `

	row := qr.QueryRowContext(ctx, query, obj.LikeableID, obj.LikeableType)
	if err := obj.Scan(row); err != nil {
		return err
	}

	return nil
}

func (obj *Like) Assign(v Like) {
	obj.ID = v.ID
	obj.LikeableID = v.LikeableID
	obj.LikeableType = v.LikeableType
	obj.CreatedAt = v.CreatedAt
	obj.UpdatedAt = v.UpdatedAt
}

func (obj Like) validate() ormerrors.ValidationError {
	return nil
}

func (obj Like) JSONAPISerialize() (string, error) {
	lines := []string{
		fmt.Sprintf("\"likeable_id\": %s", jsonapi.MustEscape(obj.LikeableID)),
		fmt.Sprintf("\"likeable_type\": %s", jsonapi.MustEscape(obj.LikeableType)),
		fmt.Sprintf("\"created_at\": %s", jsonapi.MustEscape(obj.CreatedAt)),
		fmt.Sprintf("\"updated_at\": %s", jsonapi.MustEscape(obj.UpdatedAt)),
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(lines, ", \n")), nil
}

func (obj Like) ToJSONAPIResource() (jsonapi.Resource, jsonapi.Resources) {
	includes := &jsonapi.Resources{ShouldSort: true}
	r := &jsonapi.Resource{
		ID:            jsonapi.Stringify(obj.ID),
		Type:          "like",
		Attributes:    obj,
		Relationships: jsonapi.Relationships{},
	}

	return *r, *includes
}
