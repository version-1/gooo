package schemav2

import (
	"fmt"
	"path/filepath"

	"github.com/version-1/gooo/pkg/core/generator"
	"github.com/version-1/gooo/pkg/core/schemav2/openapi"
	"github.com/version-1/gooo/pkg/core/schemav2/template"
)

type Generator struct {
	r       *openapi.RootSchema
	outputs []generator.Template
	baseURL string
	OutDir  string
}

func NewGenerator(r *openapi.RootSchema, outDir string, baseURL string) *Generator {
	return &Generator{r: r, OutDir: outDir, baseURL: baseURL}
}

func (g *Generator) Generate() error {
	schemaFile := template.SchemaFile{Schema: g.r, PackageName: "schema"}
	mainFile := template.Main{Schema: g.r}

	mainFile.Dependencies = []string{fmt.Sprintf("%s/%s", g.baseURL, filepath.Dir(schemaFile.Filename()))}

	g.outputs = append(g.outputs, schemaFile)
	g.outputs = append(g.outputs, mainFile)

	for _, tmpl := range g.outputs {
		g := generator.Generator{Dir: g.OutDir, Template: tmpl}
		if err := g.Run(); err != nil {
			return err
		}
	}

	return nil
}
