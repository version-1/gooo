package schema

import (
	"fmt"
	"strings"

	"github.com/version-1/gooo/pkg/datasource/orm/validator"
	"github.com/version-1/gooo/pkg/schema/internal/valuetype"
	gooostrings "github.com/version-1/gooo/pkg/strings"
)

type Field struct {
	Name            string
	Type            valuetype.FieldType
	TypeElementExpr string
	Tag             FieldTag
	Association     *Association
}

func (f Field) String() string {
	str := ""
	field := fmt.Sprintf("\t%s %s", f.Name, f.Type)
	str = fmt.Sprintf("%s\n", field)

	return str
}

func (f Field) ColumnName() string {
	return gooostrings.ToSnakeCase(f.Name)
}

func (f Field) TableType() string {
	v, ok := f.Type.(valuetype.FieldValueType)
	if ok {
		var opt *valuetype.FieldTableOption
		if f.Tag.TableType != "" {
			opt = &valuetype.FieldTableOption{
				Type: f.Tag.TableType,
			}
		}
		return v.TableType(opt)
	}

	return f.Type.String()
}

func (f Field) IsMutable() bool {
	return !f.Tag.Immutable && !f.Tag.Ignore
}

func (f Field) IsImmutable() bool {
	return f.Tag.Immutable && !f.Tag.Ignore
}

func (f Field) IsAssociation() bool {
	return f.Tag.Association
}

func (f Field) IsSlice() bool {
	return valuetype.MaySlice(f.Type)
}

func (f Field) IsMap() bool {
	return valuetype.MayMap(f.Type)
}

func (f Field) IsRef() bool {
	return valuetype.MayRef(f.Type)
}

func (f Field) AssociationPrimaryKey() string {
	if f.Association == nil {
		return ""
	}

	return f.Association.Schema.PrimaryKey()
}

type Validator struct {
	Fields   []string
	Validate validator.Validator
}

type Association struct {
	Slice  bool
	Schema *Schema
}

type validationKeys string

const (
	Required validationKeys = "required"
	Email    validationKeys = "email"
	Date     validationKeys = "date"
	DateTime validationKeys = "datetime"
)

type FieldTag struct {
	Raw          []string
	PrimaryKey   bool
	Immutable    bool
	Ignore       bool
	Unique       bool
	Index        bool
	DefaultValue string
	AllowNull    bool
	Association  bool
	TableType    string
	Validators   []string
}

func parseTag(tag string) FieldTag {
	if len(tag) < 2 {
		return FieldTag{}
	}
	tags := findGoooTag(tag[1 : len(tag)-1])
	options := FieldTag{
		Raw: tags,
	}
	for _, t := range tags {
		switch t {
		case "primary_key":
			options.PrimaryKey = true
		case "immutable":
			options.Immutable = true
		case "unique":
			options.Unique = true
		case "ignore":
			options.Ignore = true
		case "index":
			options.Index = true
		case "association":
			options.Association = true
		case "allow_null":
			options.AllowNull = true
		}

		if strings.HasPrefix(t, "type=") {
			segments := strings.Split(t, "=")
			if len(segments) > 1 {
				options.TableType = segments[1]
			}
		}

		if strings.HasPrefix(t, "default=") {
			segments := strings.Split(t, "=")
			if len(segments) > 1 {
				options.DefaultValue = segments[1]
			}
		}

		if strings.HasPrefix(t, "validation=") {
			segments := strings.Split(t, "=")
			if len(segments) > 1 {
				options.Validators = strings.Split(segments[1], "/")
			}
		}
	}

	return options
}

func findGoooTag(s string) []string {
	tags := strings.Split(s, " ")
	for _, t := range tags {
		parts := strings.Split(t, ":")
		if len(parts) > 1 {
			if parts[0] == "gooo" && len(parts[1]) > 2 {
				return strings.Split(parts[1][1:len(parts[1])-1], ",")
			}
		}
	}

	return []string{}
}
