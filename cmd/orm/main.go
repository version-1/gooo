package main

import (
	"os"

	exampleschema "github.com/version-1/gooo/examples/orm/schema"
)

func main() {
	args := os.Args[1:]

	dirpath := args[0]
	if err := exampleschema.Run(dirpath); err != nil {
		panic(err)
	}
}
