package schema

import "fmt"

type FieldType fmt.Stringer

type Elementer interface {
	Element() FieldType
}

type FieldValueType string

func (f FieldValueType) String() string {
	return string(f)
}

const (
	String FieldValueType = "string"
	Int    FieldValueType = "int"
	Bool   FieldValueType = "bool"
	Byte   FieldValueType = "byte"
	Time   FieldValueType = "time.Time"
	UUID   FieldValueType = "uuid.UUID"
)

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
