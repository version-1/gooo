package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	ormerrors "github.com/version-1/gooo/pkg/datasource/orm/errors"
	"github.com/version-1/gooo/pkg/datasource/query"
)

func beforeCheck[T Model](ctx context.Context, qr QueryRunner, m Model) (T, error) {
	var model T
	model, ok := m.(T)
	if !ok {
		return model, ormerrors.NewPointerModelExpectedError(m)
	}

	if err := m.Validate(); err != nil {
		return model, err
	}

	return model, nil
}

func Create[T Model](ctx context.Context, qr QueryRunner, m Model) error {
	model, err := beforeCheck[T](ctx, qr, m)
	if err != nil {
		return err
	}

	fields := m.MutableFields()
	returningKeys := m.Fields()

	columns := strings.Join(fields, ", ")
	placeholders := query.BuildPlaceholders(len(fields))
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s", m.TableName(), columns, placeholders, strings.Join(returningKeys, ", "))

	res, err := m.Scan(qr.QueryRowContext(ctx, query, m.Values()...))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ormerrors.ErrNotFound
		}
		return err
	}

	if v, ok := res.(T); ok {
		model = v
		fmt.Println(model)
	}

	return nil
}

func Update[T Model](ctx context.Context, qr QueryRunner, m Model) error {
	_, err := beforeCheck[T](ctx, qr, m)
	if err != nil {
		return err
	}

	fields := []string{}
	for i, k := range m.MutableFields() {
		// index start from 2 for id in where caluse.
		fields = append(fields, fmt.Sprintf("%s = $%d", k, i+2))
	}
	fields = append(fields, "updated_at = NOW()")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $1", m.TableName(), strings.Join(fields, ", "))

	args := append([]any{m.Identifier()}, m.Values()...)
	if _, err := qr.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	if err := Find[T](ctx, qr, m); err != nil {
		return err
	}

	return nil
}

func Delete[T Model](ctx context.Context, qr QueryRunner, m Model) error {
	_, err := beforeCheck[T](ctx, qr, m)
	if err != nil {
		return err
	}

	if err := Find[T](ctx, qr, m); err != nil {
		return err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", m.TableName())
	if _, err := qr.ExecContext(ctx, query, m.Identifier()); err != nil {
		return err
	}

	return nil
}

func Find[T Model](ctx context.Context, qr QueryRunner, m Model) error {
	var mm *T
	if _, ok := m.(T); !ok {
		return ormerrors.NewPointerModelExpectedError(m)
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1 limit 1",
		strings.Join(m.Fields(), ", "),
		m.TableName(),
	)

	res, err := m.Scan(qr.QueryRowContext(ctx, query, m.Identifier()))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ormerrors.ErrNotFound
		}

		return err
	}

	if v, ok := res.(T); ok {
		*mm = v
	}

	return nil
}

func resolveArgs(args []any) []any {
	res := make([]any, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case []uuid.UUID, []string, []int:
			res[i] = pq.Array(v)
		default:
			res[i] = v
		}
	}

	return res
}
