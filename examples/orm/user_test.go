package orm

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/version-1/gooo/pkg/datasource/logging"
	"github.com/version-1/gooo/pkg/datasource/orm"
)

func TestValidaiton(t *testing.T) {
	t.Skip("TODO:")
}

func TestUser(t *testing.T) {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	o := orm.New(db, &logging.MockLogger{}, orm.Options{QueryLog: true})

	u := NewUser()
	u.Assign(User{
		ID:       uuid.New(),
		Username: "test",
		Email:    "hoge@example.com",
	})

	ctx := context.Background()
	if err := u.Save(ctx, o); err != nil {
		t.Fatal(err)
	}

	u2 := NewUser()
	u2.Assign(
		User{
			ID: u.ID,
		})

	if err := u2.Find(ctx, o); err != nil {
		t.Fatal(err)
	}

	if u2.ID != u.ID {
		t.Fatalf("id is expected to %s, but got %s", u.ID, u2.ID)
	}

	if u2.Username != u.Username {
		t.Fatalf("username is expected to %s, but got %s", u.Username, u2.Username)
	}

	if u2.Email != u.Email {
		t.Fatalf("email is expected to %s, but got %s", u.Email, u2.Email)
	}

	if u2.CreatedAt != u.CreatedAt {
		t.Fatalf("createdAt is expected to %s, but got %s", u.CreatedAt, u2.CreatedAt)
	}

	if u2.UpdatedAt != u.UpdatedAt {
		t.Fatalf("updatedAt is expected to %s, but got %s", u.UpdatedAt, u2.UpdatedAt)
	}

	if err := u2.Destroy(ctx, o); err != nil {
		t.Fatal(err)
	}
}
