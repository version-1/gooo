package orm

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/version-1/gooo/pkg/datasource/logging"
)

func TestTransaction(t *testing.T) {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	o := New(db, &logging.MockLogger{}, Options{QueryLog: true})
	ctx := context.Background()

	if _, err := o.ExecContext(ctx, "DELETE FROM test_transaction;"); err != nil {
		t.Fatal(err)
	}

	err = o.Transaction(ctx, func(e *Executor) error {
		e.QueryRowContext(ctx, "INSERT INTO test_transaction (id) VALUES(gen_random_uuid());")
		e.QueryRowContext(ctx, "INSERT INTO test_transaction (id) VALUES(gen_random_uuid());")
		e.QueryRowContext(ctx, "INSERT INTO test_transaction (id) VALUES(gen_random_uuid());")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	var count int
	if err := o.QueryRowContext(ctx, "SELECT count(*) FROM test_transaction;").Scan(&count); err != nil {
		t.Fatal(err)
	}

	if count != 3 {
		t.Fatalf("expected 3, but got %d", count)
	}
}

func TestTransactionRollback(t *testing.T) {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	o := New(db, &logging.MockLogger{}, Options{QueryLog: true})
	ctx := context.Background()

	if _, err := o.ExecContext(ctx, "DELETE FROM test_transaction;"); err != nil {
		t.Fatal(err)
	}

	err = o.Transaction(ctx, func(e *Executor) error {
		e.QueryRowContext(ctx, "INSERT INTO test_transaction (id) VALUES(gen_random_uuid());")
		e.QueryRowContext(ctx, "INSERT INTO test_transaction (id) VALUES(gen_random_uuid());")
		e.QueryRowContext(ctx, "INSERT INTO test_transaction (id) VALUES(gen_random_uuid());")
		return errors.New("some error")
	})
	var count int
	if err := o.QueryRowContext(ctx, "SELECT count(*) FROM test_transaction;").Scan(&count); err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatalf("expected 0, but got %d", count)
	}
}
