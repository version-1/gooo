package generator

import (
	"fmt"
	"os"
)

type Generator struct {
	Dir      string
	Template Template
}

type Template interface {
	Filename() string
	Render() (string, error)
}

func (g Generator) Run() error {
	tmpl := g.Template
	filename := fmt.Sprintf("%s/%s.go", g.Dir, tmpl.Filename())
	fmt.Println("Generating: ", filename)
	s, err := g.Template.Render()
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	f.WriteString(s)

	return nil
}
