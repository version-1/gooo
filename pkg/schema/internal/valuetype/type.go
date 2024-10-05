package valuetype

import (
	"fmt"
	"go/ast"
)

type FieldType fmt.Stringer

type FieldTableOption struct {
	Type string
}

type Elementer interface {
	Element() FieldType
}

type FieldValueType string

func (f FieldValueType) String() string {
	return string(f)
}

func (f FieldValueType) TableType(option *FieldTableOption) string {
	if option != nil {
		return option.Type
	}

	switch f {
	case String:
		return "VARCHAR(255)"
	case Int:
		return "INT"
	case Bool:
		return "BOOLEAN"
	case Byte:
		return "BYTE"
	case Time:
		return "TIMESTAMP"
	case UUID:
		return "UUID"
	default:
		return f.String()
	}
}

const (
	String FieldValueType = "string"
	Int    FieldValueType = "int"
	Bool   FieldValueType = "bool"
	Byte   FieldValueType = "byte"
	Time   FieldValueType = "time.Time"
	UUID   FieldValueType = "uuid.UUID"
)

type TableFieldType string

type ref struct {
	Type FieldType
}

func (p ref) String() string {
	return fmt.Sprintf("*%s", p.Type)
}

func (p ref) Element() FieldType {
	return p.Type
}

func MayRef(f FieldType) bool {
	_, ok := f.(ref)
	return ok
}

func Ref(f FieldType) ref {
	return ref{Type: f}
}

type slice struct {
	Type FieldType
}

func (s slice) String() string {
	return fmt.Sprintf("[]%s", s.Type)
}

func (s slice) Element() FieldType {
	return s.Type
}

func MaySlice(f FieldType) bool {
	_, ok := f.(slice)
	return ok
}

func Slice(f FieldType) slice {
	return slice{Type: f}
}

type maptype struct {
	Key   FieldType
	Value FieldType
}

func (m maptype) String() string {
	return fmt.Sprintf("map[%s]%s\n", m.Key, m.Value)
}

func MayMap(f FieldType) bool {
	_, ok := f.(maptype)
	return ok
}

func Map(key, value FieldType) maptype {
	return maptype{Key: key, Value: value}
}

func convertType(s string) FieldValueType {
	switch s {
	case "string":
		return String
	case "int":
		return Int
	case "bool":
		return Bool
	case "byte":
		return Byte
	case "time.Time":
		return Time
	case "uuid.UUID":
		return UUID
	}

	return FieldValueType(s)
}

func ResolveTypeName(f ast.Expr) (FieldType, string) {
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
		tn, te := ResolveTypeName(t.X)
		typeElementExpr = te
		typeName = Ref(tn)
	case *ast.ArrayType:
		tn, te := ResolveTypeName(t.Elt)
		typeElementExpr = fmt.Sprintf("%s", tn)
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
