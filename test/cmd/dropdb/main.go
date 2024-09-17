package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

//go:embed fixtures/*.sql
var fixtures embed.FS

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	t, err := template.ParseFS(fixtures, "fixtures/*.sql")
	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	if err := t.ExecuteTemplate(&buff, "drop.sql", nil); err != nil {
		panic(err)
	}

	query := buff.String()
	fmt.Println(query)
	if _, err := db.Exec(query); err != nil {
		panic(err)
	}
}
