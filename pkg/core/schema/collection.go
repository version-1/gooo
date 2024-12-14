package schema

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/version-1/gooo/pkg/generator"
	"github.com/version-1/gooo/pkg/schema/internal/renderer"
	"github.com/version-1/gooo/pkg/util"
)

type SchemaCollection struct {
	URL     string
	Dir     string
	Package string
	Schemas []Schema
}

func (s SchemaCollection) PackageURL() string {
	url := fmt.Sprintf("%s/%s", s.URL, s.Dir)
	if strings.HasSuffix(url, "/") {
		return url[:len(url)-1]
	}

	return url
}

func (s *SchemaCollection) collect() error {
	p := NewParser()
	rootPath, err := util.LookupGomodDirPath()
	if err != nil {
		return err
	}

	path := filepath.Clean(fmt.Sprintf("%s/%s/schema.go", rootPath, s.Dir))
	list, err := p.Parse(path)
	if err != nil {
		return err
	}

	s.Schemas = list

	return nil
}

func (s SchemaCollection) schemaNames() []string {
	names := []string{}
	for _, schema := range s.Schemas {
		names = append(names, schema.Name)
	}
	return names
}

func (s SchemaCollection) Gen() error {
	if err := s.collect(); err != nil {
		return err
	}

	t := renderer.NewSharedTemplate(s.Package, s.schemaNames())
	g := generator.Generator{
		Dir:      s.Dir,
		Template: t,
	}

	if err := g.Run(); err != nil {
		return err
	}

	for _, schema := range s.Schemas {
		tmpl := renderer.SchemaTemplate{
			Basename: schema.Name,
			URL:      s.PackageURL(),
			Package:  s.Package,
			Schema:   schema,
		}

		g := generator.Generator{
			Dir:      s.Dir,
			Template: tmpl,
		}

		if err := g.Run(); err != nil {
			return err
		}
	}

	return nil
}
