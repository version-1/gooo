package jsonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Resourcer interface {
	ToJSONAPIResource() (Resource, Resources)
}

type Root[T Serializer] struct {
	Data     T
	Errors   Errors
	Meta     Serializer
	Included Resources
}

func New(data Resource, includes Resources, meta Serializer) *Root[Resource] {
	return &Root[Resource]{
		Data:     data,
		Meta:     meta,
		Included: includes,
	}
}

func NewMany(data Resources, includes Resources, meta Serializer) *Root[Resources] {
	return &Root[Resources]{
		Data:     data,
		Meta:     meta,
		Included: includes,
	}
}

func NewManyFrom[T Resourcer](list []T, meta Serializer) *Root[Resources] {
	includes := &Resources{}
	resources := &Resources{}
	for _, ele := range list {
		r, childIncludes := ele.ToJSONAPIResource()
		resources.Append(r)
		includes.Append(childIncludes.Data...)
	}

	return NewMany(*resources, *includes, meta)
}

func NewErrors(errors Errors) *Root[Nil] {
	return &Root[Nil]{
		Data:   Nil{},
		Errors: errors,
	}
}

func (j Root[T]) Serialize() (string, error) {
	fields := []string{}

	data, err := j.Data.JSONAPISerialize()
	if err != nil {
		return "", err
	}
	fields = append(fields, fmt.Sprintf("\"data\": %s", data))

	if j.Meta != nil {
		meta, err := j.Meta.JSONAPISerialize()
		if err != nil {
			return "", err
		}
		fields = append(fields, fmt.Sprintf("\"meta\": %s", meta))
	}

	errors, err := j.Errors.JSONAPISerialize()
	if err != nil {
		return "", err
	}
	fields = append(fields, fmt.Sprintf("\"errors\": %s", errors))
	included, err := j.Included.JSONAPISerialize()
	if err != nil {
		return "", err
	}
	fields = append(fields, fmt.Sprintf("\"included\": %s", included))

	s := fmt.Sprintf("{\n%s\n}", strings.Join(fields, ", \n"))

	var out bytes.Buffer
	if err := json.Indent(&out, []byte(s), "", "\t"); err != nil {
		return "", err
	}

	return out.String(), nil
}

var _ Serializer = Resource{}
var _ Serializer = Resources{}
var _ Serializer = Errors{}
var _ Serializer = Error{}
var _ Serializer = Serializers{}

type Serializer interface {
	JSONAPISerialize() (string, error)
}

type Serializers []Serializer

func (s Serializers) JSONAPISerialize() (string, error) {
	str := "["
	for _, s := range s {
		json, err := s.JSONAPISerialize()
		if err != nil {
			return "", err
		}
		str += json + ","
	}
	str += "]"
	return str, nil
}

type Resources struct {
	Data   []Resource
	keyMap map[string]bool
}

type resourceList []Resource

func (r resourceList) Len() int { return len(r) }
func (r resourceList) Less(i, j int) bool {
	if r[i].Type != r[j].Type {
		return r[i].Type < r[j].Type
	}

	return r[i].ID < r[j].ID
}
func (r resourceList) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func (j *Resources) Append(r ...Resource) {
	if j.keyMap == nil {
		j.keyMap = make(map[string]bool)
	}

	for _, res := range r {
		key := fmt.Sprintf("%s:%s", res.Type, res.ID)
		if _, ok := j.keyMap[key]; !ok {
			j.Data = append(j.Data, res)
			j.keyMap[key] = true
		}
	}

	sort.Sort(resourceList(j.Data))
}

func (j Resources) JSONAPISerialize() (string, error) {
	lines := []string{}
	for _, r := range j.Data {
		json, err := r.JSONAPISerialize()
		if err != nil {
			return "", err
		}
		lines = append(lines, json)
	}

	str := "["
	str += strings.Join(lines, ", \n")
	str += "]"
	return str, nil
}

type Resource struct {
	ID            string
	Type          string
	Attributes    Serializer
	Relationships Relationships
}

func (j Resource) JSONAPISerialize() (string, error) {
	attrs, err := j.Attributes.JSONAPISerialize()
	if err != nil {
		return "", err
	}

	r, err := j.Relationships.JSONAPISerialize()
	if err != nil {
		return "", err
	}

	return `{
		"id": ` + j.ID + `,
		"type": "` + j.Type + `",
		"attributes": ` + attrs + `,
		"relationships": ` + r + `
	}`, nil
}

type Relationships map[string]Serializer

func (j Relationships) JSONAPISerialize() (string, error) {
	lines := []string{}
	for k, r := range j {
		json, err := r.JSONAPISerialize()
		if err != nil {
			return "", err
		}
		lines = append(lines, "\""+k+"\": "+json)
	}
	str := "{"
	str += strings.Join(lines, ", \n")
	str += "}"
	return str, nil
}

type RelationshipHasMany struct {
	Data []ResourceIdentifier
}

func (j RelationshipHasMany) JSONAPISerialize() (string, error) {
	lines := []string{}
	for i := range j.Data {
		json, err := j.Data[i].JSONAPISerialize()
		if err != nil {
			return "", err
		}

		lines = append(lines, json)
	}
	str := "["
	str += strings.Join(lines, ", \n")
	str += "]"

	return `{
		"data": ` + str + `
	}`, nil
}

type Relationship struct {
	Data ResourceIdentifier
}

func (j Relationship) JSONAPISerialize() (string, error) {
	json, err := j.Data.JSONAPISerialize()
	if err != nil {
		return "", err
	}

	return `{
		"data": ` + json + `
	}`, nil
}

type ResourceIdentifier struct {
	ID   string
	Type string
}

func (j ResourceIdentifier) JSONAPISerialize() (string, error) {
	return `{
		"id": ` + j.ID + `,
		"type": "` + j.Type + `"
	}`, nil
}

type Errors []Error

func (j Errors) JSONAPISerialize() (string, error) {
	str := "["
	for _, e := range j {
		json, err := e.JSONAPISerialize()
		if err != nil {
			return "", err
		}
		str += json + ","
	}
	str += "]"

	return str, nil
}

type Error struct {
	ID     string
	Status int
	Code   string
	Title  string
	Detail string
}

type Nil struct{}

func (n Nil) JSONAPISerialize() (string, error) {
	return "null", nil
}

func (j Error) JSONAPISerialize() (string, error) {
	fields := []string{
		fmt.Sprintf("\"id\": %s", Stringify(j.ID)),
		fmt.Sprintf("\"status\": %s", Stringify(j.Status)),
		fmt.Sprintf("\"code\": %s", Stringify(j.Code)),
		fmt.Sprintf("\"title\": %s", Stringify(j.Title)),
		fmt.Sprintf("\"detail\": %s", Stringify(j.Detail)),
	}

	return fmt.Sprintf("{\n%s\n}", strings.Join(fields, ", \n")), nil
}
