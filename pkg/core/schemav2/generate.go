package schemav2

import (
	"go/format"

	"github.com/version-1/gooo/pkg/core/generator"
	"golang.org/x/tools/imports"
)

type Generator struct {
	r       *RootSchema
	outputs []generator.Template
	OutDir  string
}

func NewGenerator(r *RootSchema, outDir string) *Generator {
	return &Generator{r: r, OutDir: outDir}
}

func (g *Generator) Generate() error {
	g.outputs = append(g.outputs, Main{
		Routes: "// ここにルーティングが入ります",
	})
	for _, tmpl := range g.outputs {
		g := generator.Generator{Dir: g.OutDir, Template: tmpl}
		if err := g.Run(); err != nil {
			return err
		}
	}

	return nil
}

func pretify(filename, s string) ([]byte, error) {
	formatted, err := format.Source([]byte(s))
	if err != nil {
		return []byte{}, err
	}

	processed, err := imports.Process(filename, formatted, nil)
	if err != nil {
		return formatted, err
	}

	return processed, nil
}
