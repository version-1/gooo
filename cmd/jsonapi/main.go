package jsonapi

import "os"

func main() {
	args := os.Args[1:]

	dirpath := args[0]
	if err := exampleschema.Run(dirpath); err != nil {
		panic(err)
	}

}
