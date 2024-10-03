package schema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser_Parse(t *testing.T) {
	p := NewParser()
	list, err := p.Parse("./internal/fixtures/schema.go")
	if err != nil {
		t.Fatal(err)
	}

	expect := []Schema{
		{
			Name:      "User",
			TableName: "users",
			Fields: []Field{
				{
					Name:            "ID",
					Type:            FieldType(Int),
					TypeElementExpr: "int",
					Tag: FieldTag{
						Raw:        []string{"primary_key", "immutable"},
						PrimaryKey: true,
						Immutable:  true,
					},
				},
				{
					Name:            "Username",
					Type:            FieldType(String),
					TypeElementExpr: "string",
					Tag: FieldTag{
						Raw:    []string{"unique"},
						Unique: true,
					},
				},
				{
					Name:            "Email",
					Type:            FieldType(String),
					TypeElementExpr: "string",
					Tag: FieldTag{
						Raw: []string{},
					},
				},
				{
					Name:            "RefreshToken",
					Type:            FieldType(String),
					TypeElementExpr: "string",
					Tag: FieldTag{
						Raw: []string{},
					},
				},
				{
					Name:            "Timezone",
					Type:            FieldType(String),
					TypeElementExpr: "string",
					Tag: FieldTag{
						Raw: []string{},
					},
				},
				{
					Name:            "TimeDiff",
					Type:            FieldType(Int),
					TypeElementExpr: "int",
					Tag: FieldTag{
						Raw: []string{},
					},
				},
				{
					Name:            "CreatedAt",
					Type:            FieldType(Time),
					TypeElementExpr: "time.Time",
					Tag: FieldTag{
						Raw:       []string{"immutable"},
						Immutable: true,
					},
				},
				{
					Name:            "UpdatedAt",
					Type:            FieldType(Time),
					TypeElementExpr: "time.Time",
					Tag: FieldTag{
						Raw:       []string{"immutable"},
						Immutable: true,
					},
				},
				{
					Name:            "Profile",
					Type:            Ref(FieldValueType("Profile")),
					TypeElementExpr: "Profile",
					Tag: FieldTag{
						Raw:         []string{"association"},
						Association: true,
					},
				},
				{
					Name:            "Posts",
					Type:            Slice(FieldValueType("Post")),
					TypeElementExpr: "[]Post",
					Tag: FieldTag{
						Raw:         []string{"association"},
						Association: true,
					},
				},
			},
		},
	}

	actual := list[0:1]
	if diff := cmp.Diff(expect, actual); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
