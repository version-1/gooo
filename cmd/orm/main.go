package main

import (
	"os"

	"github.com/version-1/gooo/examples/schema"
)

func main() {
	args := os.Args[1:]

	dirpath := args[0]
	if err := schema.Run(dirpath); err != nil {
		panic(err)
	}
}
