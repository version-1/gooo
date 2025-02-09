package jsonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	goooerrors "github.com/version-1/gooo/pkg/toolkit/errors"
	"github.com/version-1/gooo/pkg/toolkit/logger"
)

type Resourcer interface {
	ToJSONAPIResource() (data Resource, included Resources)
}

type Root[T Serializer] struct {
	Data     T
	Errors   Errors
	Meta     Serializer
	Included Resources
}

func New(data Resource, includes Resources, meta Serializer) (*Root[Resource], error) {
	if data.ID == "" {
		return nil, goooerrors.Errorf("ID is required.")
	}

	return &Root[Resource]{
		Data:     data,
		Meta:     meta,
		Included: includes,
	}, nil
}

func newMany(data Resources, includes Resources, meta Serializer) *Root[Resources] {
	return &Root[Resources]{
		Data:     data,
		Meta:     meta,
		Included: includes,
	}
}

func NewManyFrom[T Resourcer](list []T, meta Serializer) (*Root[Resources], error) {
	includes := &Resources{
		ShouldSort: true,
	}
	resources := &Resources{}
	for index, ele := range list {
		r, childIncludes := ele.ToJSONAPIResource()
		if r.ID == "" {
			return nil, goooerrors.Errorf("ID is required. index: %d", index)
		}
		resources.Append(r)
		includes.Append(childIncludes.Data...)
	}

	return newMany(*resources, *includes, meta), nil
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
		return "", goooerrors.Wrap(err)
	}
	fields = append(fields, fmt.Sprintf("\"data\": %s", data))

	if j.Meta != nil {
		meta, err := j.Meta.JSONAPISerialize()
		if err != nil {
			return "", goooerrors.Wrap(err)
		}
		fields = append(fields, fmt.Sprintf("\"meta\": %s", meta))
	}

	errors, err := j.Errors.JSONAPISerialize()
	if err != nil {
		return "", goooerrors.Wrap(err)
	}

	if errors != "[]" {
		fields = append(fields, fmt.Sprintf("\"errors\": %s", errors))
	}

	included, err := j.Included.JSONAPISerialize()
	if err != nil {
		return "", goooerrors.Wrap(err)
	}

	if included != "[]" {
		fields = append(fields, fmt.Sprintf("\"included\": %s", included))
	}

	s := fmt.Sprintf("{\n%s\n}", strings.Join(fields, ", \n"))

	var out bytes.Buffer
	if err := json.Compact(&out, []byte(s)); err != nil {
		logger.DefaultLogger.Errorf("pkg/presenter/jsonapi: got error on compact json. %s")
		logger.DefaultLogger.Errorf("%s\n", s)
		return "", goooerrors.Wrap(err)
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
	for i, serializer := range s {
		json, err := serializer.JSONAPISerialize()
		if err != nil {
			return "", goooerrors.Wrap(err)
		}

		comma := ""
		if i != len(s)-1 {
			comma = ","
		}
		str += json + comma
	}
	str += "]"
	return str, nil
}

type Attributes[T any] struct {
	v T
}

func NewAttributes[T any](v T) Attributes[T] {
	return Attributes[T]{v}
}

func (a Attributes[T]) JSONAPISerialize() (string, error) {
	if b, err := json.Marshal(a.v); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
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

type Resources struct {
	Data       []Resource
	keyMap     map[string]bool
	ShouldSort bool
}

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

	if j.ShouldSort {
		sort.Sort(resourceList(j.Data))
	}
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
		return "", goooerrors.Wrap(err)
	}

	r, err := j.Relationships.JSONAPISerialize()
	if err != nil {
		return "", goooerrors.Wrap(err)
	}

	list := []string{}

	id, err := Escape(j.ID)
	if err != nil {
		return "", goooerrors.Wrap(err)
	}

	t, err := Escape(j.Type)
	if err != nil {
		return "", goooerrors.Wrap(err)
	}
	list = append(list, fmt.Sprintf(`"id": %s`, id))
	list = append(list, fmt.Sprintf(`"type": %s`, t))
	list = append(list, fmt.Sprintf(`"attributes": %s`, attrs))
	if r != "{}" {
		list = append(list, fmt.Sprintf(`"relationships": %s`, r))
	}

	return "{\n" + strings.Join(list, ", \n") + "\n}", nil
}

func roundQuotes(s string) string {
	if len(s) < 1 {
		return s
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}

	return s
}

type Relationships map[string]Serializer

func (j Relationships) JSONAPISerialize() (string, error) {
	lines := []string{}
	for k, r := range j {
		json, err := r.JSONAPISerialize()
		if err != nil {
			return "", goooerrors.Wrap(err)
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
			return "", goooerrors.Wrap(err)
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
		return "", goooerrors.Wrap(err)
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
	id, err := Escape(j.ID)
	if err != nil {
		return "", goooerrors.Wrap(err)
	}

	t, err := Escape(j.Type)
	if err != nil {
		return "", goooerrors.Wrap(err)
	}
	return `{
		"id": ` + id + `,
		"type": ` + t + `
	}`, nil
}

type Nil struct{}

func (n Nil) JSONAPISerialize() (string, error) {
	return "null", nil
}

func HasOne(r *Resource, includes *Resources, ele Resourcer, id any, typeName string) {
	if ele == nil {
		return
	}

	r.Relationships[typeName] = Relationship{
		Data: ResourceIdentifier{
			ID:   Stringify(id),
			Type: typeName,
		},
	}

	resource, childIncludes := ele.ToJSONAPIResource()
	includes.Append(resource)
	includes.Append(childIncludes.Data...)
}

func HasMany(r *Resource, includes *Resources, elements []Resourcer, typeName string, cb func(ri *ResourceIdentifier, index int)) {
	relationships := RelationshipHasMany{}
	for i, ele := range elements {
		ri := ResourceIdentifier{
			Type: typeName,
		}
		cb(&ri, i)
		relationships.Data = append(
			relationships.Data,
			ri,
		)

		resource, childIncludes := ele.ToJSONAPIResource()
		includes.Append(resource)
		includes.Append(childIncludes.Data...)
	}

	if len(relationships.Data) > 0 {
		r.Relationships[typeName] = relationships
	}
}
