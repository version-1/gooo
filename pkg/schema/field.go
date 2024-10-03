package schema

import (
	"strings"
)

type validationKeys string

const (
	Required validationKeys = "required"
	Email    validationKeys = "email"
	Date     validationKeys = "date"
	DateTime validationKeys = "datetime"
)

type FieldTag struct {
	Raw         []string
	PrimaryKey  bool
	Immutable   bool
	Ignore      bool
	Unique      bool
	Index       bool
	Association bool
	TableType   string
	Validators  []string
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
		}

		if strings.HasPrefix(t, "type=") {
			segments := strings.Split(t, "=")
			if len(segments) > 1 {
				options.TableType = segments[1]
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
