package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/version-1/gooo/pkg/command/migration"
	"github.com/version-1/gooo/pkg/command/migration/runner"
)

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	m, err := runner.NewYaml(db, os.Getenv("MIGRATION_PATH"))
	if err != nil {
		panic(err)
	}

	c, err := migration.NewWith(db, m, nil)
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 1 {
		fmt.Println("command is required. [up|down|create|drop|generate]")
		os.Exit(1)
		return
	}

	cmd := os.Args[1]

	args := []string{}
	if len(os.Args) >= 3 {
		args = os.Args[2:]
	}

	if err := c.Exec(ctx, cmd, args...); err != nil {
		panic(err)
	}
}
