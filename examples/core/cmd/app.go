package main

import (
	schema "github.com/version-1/gooo/pkg/core/schemav2"
)

func main() {
	s, err := schema.New("./examples/bare/internal/swagger/swagger.yml")
	if err != nil {
		panic(err)
	}

	g := schema.NewGenerator(s, "./examples/core/generated")

	if err := g.Generate(); err != nil {
		panic(err)
	}
}
