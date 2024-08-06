package schema

import "fmt"

type FieldType fmt.Stringer

type FieldValueType string

func (f FieldValueType) String() string {
	switch f {
	case UUID:
		return "uuid.UUID"
	}
	return string(f)
}

const (
	UUID   FieldValueType = "uuid"
	Int    FieldValueType = "int"
	String FieldValueType = "string"
	Time   FieldValueType = "time.Time"
)

type Ptr struct {
	Type FieldType
}

func (p Ptr) String() string {
	return fmt.Sprintf("*%s", p.Type)
}

func Ref(f FieldType) Ptr {
	return Ptr{Type: f}
}
