package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/version-1/gooo/pkg/errors"
	"github.com/version-1/gooo/pkg/presenter/jsonapi"
	goootesting "github.com/version-1/gooo/pkg/testing"
)

type dummy struct {
	ID     string    `json:"-"`
	String string    `json:"string"`
	Number int       `json:"number"`
	Bool   bool      `json:"bool"`
	Time   time.Time `json:"time"`
}

func (d dummy) ToJSONAPIResource() (jsonapi.Resource, jsonapi.Resources) {
	return jsonapi.Resource{
		ID:         d.ID,
		Type:       "dummy",
		Attributes: jsonapi.NewAttributes(d),
	}, jsonapi.Resources{}
}

type meta struct {
	Key string `json:"key"`
}

func (m meta) JSONAPISerialize() (string, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

func TestJSONAPIContentType(t *testing.T) {
	a := JSONAPI{}
	expect := "application/json"
	if a.ContentType() != expect {
		t.Errorf("Expected content type to be %s, got %s", expect, a.ContentType())
	}
}

func TestJSONAPIRender(t *testing.T) {
	a := JSONAPI{}
	id1 := uuid.MustParse("325fe993-420a-4e53-8687-1760f34e0697").String()
	id2 := uuid.MustParse("e3a341b2-0400-4e80-97b9-b1aa0119018b").String()
	id3 := uuid.MustParse("f513710d-a158-4cdb-914f-bb8aa11bd675").String()
	now := time.Now()

	test := goootesting.NewTable([]goootesting.Record[[]byte, []byte]{
		{
			Name: "Render with jsonapi.Resourcer",
			Subject: func(t *testing.T) ([]byte, error) {
				s, err := a.Render(dummy{
					ID:     id1,
					String: "string",
					Number: 1,
					Bool:   true,
					Time:   now,
				})
				if err != nil {
					return []byte{}, err
				}

				buffer := &bytes.Buffer{}
				err = json.Compact(buffer, s)
				return buffer.Bytes(), err
			},
			Expect: func(t *testing.T) ([]byte, error) {
				s := fmt.Sprintf(`{ "data": { "id": "%s", "type": "dummy", "attributes": { "string": "string", "number": 1, "bool": true, "time": "%s" } } }`, id1, now.Format(time.RFC3339Nano))
				buffer := &bytes.Buffer{}
				err := json.Compact(buffer, []byte(s))
				return buffer.Bytes(), err
			},
			Assert: func(t *testing.T, r *goootesting.Record[[]byte, []byte]) bool {
				e, err := r.Expect(t)
				s, serr := r.Subject(t)

				if !reflect.DeepEqual(e, s) {
					t.Errorf("Expected %s, got %s", e, s)
					return false
				}

				if serr != nil && err.Error() != serr.Error() {
					t.Errorf("Expected %v, got %v", err, serr)
					return false
				}
				return true
			},
		},
		{
			Name: "Render with []jsonapi.Resourcer",
			Subject: func(t *testing.T) ([]byte, error) {
				list := []jsonapi.Resourcer{
					dummy{
						ID:     id1,
						String: "string",
						Number: 1,
						Bool:   true,
						Time:   now,
					},
					dummy{
						ID:     id2,
						String: "string",
						Number: 2,
						Bool:   true,
						Time:   now,
					},
					dummy{
						ID:     id3,
						String: "string",
						Number: 3,
						Bool:   true,
						Time:   now,
					},
				}
				s, err := a.Render(list)
				if err != nil {
					return []byte{}, err
				}

				buffer := &bytes.Buffer{}
				err = json.Compact(buffer, s)
				return buffer.Bytes(), err
			},
			Expect: func(t *testing.T) ([]byte, error) {
				s := fmt.Sprintf(`{
					"data": [
						{ "id": "%s", "type": "dummy", "attributes": { "string": "string", "number": 1, "bool": true, "time": "%s" } },
						{ "id": "%s", "type": "dummy", "attributes": { "string": "string", "number": 2, "bool": true, "time": "%s" } },
						{ "id": "%s", "type": "dummy", "attributes": { "string": "string", "number": 3, "bool": true, "time": "%s" } }
					]
				}`,
					id1,
					now.Format(time.RFC3339Nano),
					id2,
					now.Format(time.RFC3339Nano),
					id3,
					now.Format(time.RFC3339Nano),
				)

				buffer := &bytes.Buffer{}
				err := json.Compact(buffer, []byte(s))
				return buffer.Bytes(), err
			},
			Assert: func(t *testing.T, r *goootesting.Record[[]byte, []byte]) bool {
				e, err := r.Expect(t)
				s, serr := r.Subject(t)

				if !reflect.DeepEqual(e, s) {
					t.Errorf("Expected %s, got %s", e, s)
					return false
				}

				if serr != nil && err.Error() != serr.Error() {
					t.Errorf("Expected %v, got %v", err, serr)
					return false
				}
				return true
			},
		},
		{
			Name: "Render with []jsonapi.Resourcer and meta",
			Subject: func(t *testing.T) ([]byte, error) {
				list := []jsonapi.Resourcer{
					dummy{
						ID:     id2,
						String: "string",
						Number: 1,
						Bool:   true,
						Time:   now,
					},
					dummy{
						ID:     id1,
						String: "string",
						Number: 2,
						Bool:   true,
						Time:   now,
					},
					dummy{
						ID:     id3,
						String: "string",
						Number: 3,
						Bool:   true,
						Time:   now,
					},
				}

				option := JSONAPIOption{
					Meta: meta{
						Key: "value",
					},
				}
				s, err := a.Render(list, option)
				if err != nil {
					return []byte{}, err
				}

				buffer := &bytes.Buffer{}
				err = json.Compact(buffer, s)
				return buffer.Bytes(), err
			},
			Expect: func(t *testing.T) ([]byte, error) {
				s := fmt.Sprintf(`{
					"data": [
						{ "id": "%s", "type": "dummy", "attributes": { "string": "string", "number": 1, "bool": true, "time": "%s" } },
						{ "id": "%s", "type": "dummy", "attributes": { "string": "string", "number": 2, "bool": true, "time": "%s" } },
						{ "id": "%s", "type": "dummy", "attributes": { "string": "string", "number": 3, "bool": true, "time": "%s" } }
					],
					"meta": { "key": "value" }
				}`,
					id2,
					now.Format(time.RFC3339Nano),
					id1,
					now.Format(time.RFC3339Nano),
					id3,
					now.Format(time.RFC3339Nano),
				)

				buffer := &bytes.Buffer{}
				err := json.Compact(buffer, []byte(s))
				return buffer.Bytes(), err
			},
			Assert: func(t *testing.T, r *goootesting.Record[[]byte, []byte]) bool {
				e, err := r.Expect(t)
				s, serr := r.Subject(t)

				if !reflect.DeepEqual(e, s) {
					t.Errorf("Expected %s, got %s", e, s)
					return false
				}

				if serr != nil && err.Error() != serr.Error() {
					t.Errorf("Expected %v, got %v", err, serr)
					return false
				}
				return true
			},
		},
		{
			Name: "Render with invalid type",
			Subject: func(t *testing.T) ([]byte, error) {
				return a.Render("hoge")
			},
			Expect: func(t *testing.T) ([]byte, error) {
				return []byte{}, errors.Wrap(JSONAPIInvalidTypeError{"hoge"})
			},
			Assert: func(t *testing.T, r *goootesting.Record[[]byte, []byte]) bool {
				e, err := r.Expect(t)
				s, serr := r.Subject(t)

				if !reflect.DeepEqual(e, s) {
					t.Errorf("Expected %v, got %v", e, err)
					return false
				}

				if err.Error() != serr.Error() {
					t.Errorf("Expected %v, got %v", err, serr)
					return false
				}
				return true
			},
		},
	})

	test.Run(t)
}
