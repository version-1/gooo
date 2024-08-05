package orm

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
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
		return errors.New("primaryKey is required")
	}

	query := "DELETE FROM users WHERE id = $1"
	if _, err := qr.ExecContext(ctx, query, obj.ID); err != nil {
		return err
	}

	return nil
}

func (obj *User) Find(ctx context.Context, qr queryer) error {
	if obj.ID == uuid.Nil {
		return errors.New("primaryKey is required")
	}

	query := "SELECT id, username, bio, email, created_at, updated_at FROM users WHERE id = $1"
	row := qr.QueryRowContext(ctx, query, obj.ID)

	if err := obj.Scan(row); err != nil {
		return err
	}

	return nil
}

func (obj *User) Save(ctx context.Context, qr queryer) error {
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
