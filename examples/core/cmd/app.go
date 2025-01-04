package main

import (
	"fmt"

	schema "github.com/version-1/gooo/pkg/core/schemav2"
)

func main() {
	s, err := schema.New("./examples/bare/internal/swagger/swagger.yml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", s)
}
