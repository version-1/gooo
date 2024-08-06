package orm

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	goooerrors "github.com/version-1/gooo/pkg/datasource/orm/errors"
	"github.com/version-1/gooo/pkg/datasource/schema"
)

type User struct {
	schema.Schema
	ID        uuid.UUID
	Username  string
	Bio       *string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
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
