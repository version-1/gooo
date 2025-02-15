package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/version-1/gooo/pkg/errors"
	"github.com/version-1/gooo/pkg/util"
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
	relativePath := filepath.Clean(fmt.Sprintf("%s/%s.go", g.Dir, tmpl.Filename()))
	rootPath, err := util.LookupGomodDirPath()
	if err != nil {
		return err
	}
	filename := filepath.Clean(fmt.Sprintf("%s/%s", rootPath, relativePath))
	fmt.Println("Generating: ", relativePath)
	s, err := g.Template.Render()
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err)
	}

	defer f.Close()

	f.WriteString(s)

	return nil
}
