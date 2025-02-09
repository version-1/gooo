package main

import (
	"fmt"

	schema "github.com/version-1/gooo/pkg/core/schema"
	"github.com/version-1/gooo/pkg/core/schema/openapi/v3_0_0"
)

func main() {
	s, err := v3_0_0.New("./examples/bare/internal/swagger/swagger.yml")
	if err != nil {
		panic(err)
	}

	g := schema.NewGenerator(s, "./examples/core/generated", "github.com/version-1/gooo/examples/core/generated")

	if err := g.Generate(); err != nil {
		fmt.Printf("Error: %+v\n", err)
		panic(err)
	}
}
