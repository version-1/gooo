package schema

import "fmt"

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
