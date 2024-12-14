package schema

import (
	"go/ast"
	"go/token"
	"os"

	goparser "go/parser"

	"github.com/version-1/gooo/pkg/errors"
	"github.com/version-1/gooo/pkg/schema/internal/valuetype"
	"github.com/version-1/gooo/pkg/strings"
)

type parser struct{}

func NewParser() *parser {
	return &parser{}
}

func (p parser) Parse(path string) ([]Schema, error) {
	list := []Schema{}
	fset := token.NewFileSet()
	src, err := os.ReadFile(path)
	if err != nil {
		return list, errors.Wrap(err)
	}

	node, err := goparser.ParseFile(fset, "", src, goparser.ParseComments)
	if err != nil {
		return list, errors.Wrap(err)
	}

	m := map[string]*Schema{}
	ast.Inspect(node, func(n ast.Node) bool {
		if t, ok := n.(*ast.TypeSpec); ok {
			name := t.Name.Name
			if len(list) > 0 {
				m[list[len(list)-1].Name] = &list[len(list)-1]
			}
			list = append(list, Schema{
				Name:      name,
				TableName: strings.ToPlural(name),
			})
		}

		if field, ok := n.(*ast.Field); ok {
			if field.Tag != nil {
				typeName, typeElementExpr := valuetype.ResolveTypeName(field.Type)
				list[len(list)-1].AddFields(Field{
					Name:            field.Names[0].Name,
					Type:            typeName,
					TypeElementExpr: typeElementExpr,
					Tag:             parseTag(field.Tag.Value),
				})
			}
		}
		return true
	})

	m[list[len(list)-1].Name] = &list[len(list)-1]

	for i := range list {
		for j := range list[i].Fields {
			f := list[i].Fields[j]
			if f.IsAssociation() {
				schema, ok := m[f.TypeElementExpr]
				if !ok {
					return list, errors.Errorf("schema %s not found on association", f.TypeElementExpr)
				}

				list[i].Fields[j].Association = &Association{
					Schema: schema,
					Slice:  f.IsSlice(),
				}
			}
		}
	}

	return list, nil
}
