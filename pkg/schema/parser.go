package schema

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"

	goparser "go/parser"

	"github.com/version-1/gooo/pkg/strings"
)

type parser struct {
}

func NewParser() *parser {
	return &parser{}
}

func (p parser) Parse(path string) ([]Schema, error) {
	list := []Schema{}
	fset := token.NewFileSet()
	src, err := os.ReadFile(path)
	if err != nil {
		return list, err
	}

	node, err := goparser.ParseFile(fset, "", src, goparser.ParseComments)
	if err != nil {
		return list, err
	}

	m := map[string]Schema{}
	ast.Inspect(node, func(n ast.Node) bool {
		if t, ok := n.(*ast.TypeSpec); ok {
			name := t.Name.Name
			if len(list) > 0 {
				m[list[len(list)-1].Name] = list[len(list)-1]
			}
			list = append(list, Schema{
				Name:      name,
				TableName: strings.ToPlural(name),
			})
		}

		if field, ok := n.(*ast.Field); ok {
			if field.Tag != nil {
				typeName, typeElementExpr := resolveTypeName(field.Type)
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

	for _, s := range list {
		for _, f := range s.Fields {
			if f.IsAssociation() {
				f.Association = &Association{
					Schema: m[f.Type.String()],
					Slice:  f.IsSlice(),
				}
			}
		}
	}

	return list, nil
}

func resolveTypeName(f ast.Expr) (FieldType, string) {
	var typeName FieldType
	var typeElementExpr string
	switch t := f.(type) {
	case *ast.Ident:
		typeElementExpr = t.Name
		typeName = convertType(typeElementExpr)
	case *ast.SelectorExpr:
		typeElementExpr = fmt.Sprintf("%s.%s", t.X, t.Sel)
		typeName = convertType(typeElementExpr)
	case *ast.StarExpr:
		tn, te := resolveTypeName(t.X)
		typeElementExpr = te
		typeName = Ref(tn)
	case *ast.ArrayType:
		tn, te := resolveTypeName(t.Elt)
		typeElementExpr = fmt.Sprintf("[]%s", tn)
		typeName = Slice(convertType(te))
	case *ast.MapType:
		typeName = Map(
			convertType(fmt.Sprintf("%s", t.Key)),
			convertType(fmt.Sprintf("%s", t.Value)),
		)
		typeElementExpr = typeName.String()
	}

	return typeName, typeElementExpr
}
