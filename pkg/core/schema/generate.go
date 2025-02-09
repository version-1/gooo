package schema

import (
	"fmt"
	"path/filepath"

	"github.com/version-1/gooo/pkg/core/generator"
	"github.com/version-1/gooo/pkg/core/schema/openapi/v3_0_0"
	"github.com/version-1/gooo/pkg/core/schema/template"
)

type Generator struct {
	r       *v3_0_0.RootSchema
	outputs []generator.Template
	baseURL string
	OutDir  string
}

func NewGenerator(r *v3_0_0.RootSchema, outDir string, baseURL string) *Generator {
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
