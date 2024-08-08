package jsonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Resourcer interface {
	ToJSONAPIResource() (Resource, Resources)
}

type Root[T Serializer] struct {
	Data     T
	Errors   ErrorSerializers
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

func NewErrors(errors []ErrorSerializer, meta Serializer) *Root[Nil] {
	return &Root[Nil]{
		Data:   Nil{},
		Errors: ErrorSerializers(errors),
		Meta:   meta,
	}
}

func (j Root[T]) Serialize() (string, error) {
	fields := []string{}

	data, err := j.Data.JSONAPISerialize()
	if err != nil {
		return "", err
	}

	d := strings.TrimSpace(data)
	if !isEmptyJSON(d) {
		fields = append(fields, fmt.Sprintf("\"data\": %s", d))
	}

	if j.Meta != nil {
		meta, err := j.Meta.JSONAPISerialize()
		if err != nil {
			return "", err
		}

		d := strings.TrimSpace(meta)
		fields = append(fields, fmt.Sprintf("\"meta\": %s", d))
	}

	if len(j.Errors) > 0 {
		errors, err := j.Errors.JSONAPIErrorSerialize()
		if err != nil {
			return "", err
		}

		d = strings.TrimSpace(errors)
		fields = append(fields, fmt.Sprintf("\"errors\": %s", d))
	}

	if len(j.Included.Data) > 0 {
		included, err := j.Included.JSONAPISerialize()
		if err != nil {
			return "", err
		}

		d = strings.TrimSpace(included)
		fields = append(fields, fmt.Sprintf("\"included\": %s", d))
	}

	s := fmt.Sprintf("{\n%s\n}", strings.Join(fields, ",\n"))

	var out bytes.Buffer
	if err := json.Indent(&out, []byte(s), "", "  "); err != nil {
		return "", err
	}

	return out.String(), nil
}

var _ Serializer = Resource{}
var _ Serializer = Resources{}
var _ Serializer = Error{}
var _ Serializer = Serializers{}
var _ Serializer = ErrorSerializers{}

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

type ErrorSerializer interface {
	JSONAPIErrorSerialize() (string, error)
}

type ErrorSerializers []ErrorSerializer

func (j ErrorSerializers) JSONAPISerialize() (string, error) {
	elements := []string{}
	for _, e := range j {
		json, err := e.JSONAPIErrorSerialize()
		if err != nil {
			return "", err
		}
		elements = append(elements, json)
	}

	return renderList(elements), nil
}

func (j ErrorSerializers) JSONAPIErrorSerialize() (string, error) {
	return j.JSONAPISerialize()
}

type Error struct {
	ID     string
	Status int
	Code   string
	Title  string
	Detail string
}

func (j Error) JSONAPIErrorSerialize() (string, error) {
	return ObjectBuilder{}.Append(
		NewField("id", j.ID),
		NewField("title", j.Title),
		NewField("code", j.Code),
		NewField("status", strconv.Itoa(j.Status)),
		NewField("detail", j.Detail),
	).String(), nil
}

func (j Error) JSONAPISerialize() (string, error) {
	return j.JSONAPIErrorSerialize()
}

var _ error = ErrorWrapper{}
var _ ErrorSerializer = ErrorWrapper{}

type ErrorWrapper struct {
	Err    error
	ID     string
	Title  string
	Code   string
	Status int
}

func (e ErrorWrapper) JSONAPIErrorSerialize() (string, error) {
	return Error{
		ID:     e.ID,
		Title:  e.Title,
		Code:   e.Code,
		Status: e.Status,
		Detail: e.Detail(),
	}.JSONAPIErrorSerialize()
}

func (e ErrorWrapper) Error() string {
	return e.Err.Error()
}

func (e ErrorWrapper) Detail() string {
	return e.Error()
}

type Nil struct{}

func (n Nil) JSONAPISerialize() (string, error) {
	return "null", nil
}
