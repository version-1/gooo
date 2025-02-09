package payload

import (
	"testing"
)

type ConfigKey string

const (
	PORT        ConfigKey = "PORT"
	DATABAE_URL ConfigKey = "DATABASE_URL"
)

func TestLoad(t *testing.T) {
	t.Skip("skipping test in CI")
	loader := NewEnvfileLoader[ConfigKey]("./fixtures/.env.test")
	m, err := loader.Load()
	if err != nil {
		t.Fatal(err)
	}

	i := 0
	for k, v := range *m {
		if k == "PORT" && v != "3000" {
			t.Fatalf("expected %s, got %s", "3000", v)
		}

		if k == "DATABASE_URL" && v != "postgres://postgres:password@localhost:5432/test?sslmode=disable" {
			t.Fatalf("expected %s, got %s", "postgres://postgres:password@localhost:5432/test?sslmode=disabled", v)
		}

		i++
	}

	if i != 3 {
		t.Fatalf("expected %d, got %d", 3, i)
	}
}
