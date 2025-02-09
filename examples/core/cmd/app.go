package main

import (
	"fmt"

	schema "github.com/version-1/gooo/pkg/core/schema"
	"github.com/version-1/gooo/pkg/core/schema/openapi"
)

func main() {
	s, err := openapi.New("./examples/bare/internal/swagger/swagger.yml")
	if err != nil {
		panic(err)
	}

	g := schema.NewGenerator(s, "./examples/core/generated", "github.com/version-1/gooo/examples/core")

	if err := g.Generate(); err != nil {
		fmt.Printf("Error: %+v\n", err)
		panic(err)
	}
}
