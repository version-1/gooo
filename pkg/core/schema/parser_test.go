package schema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/version-1/gooo/pkg/schema/internal/valuetype"
)

func TestParser_Parse(t *testing.T) {
	p := NewParser()
	list, err := p.Parse("./internal/schema/schema.go")
	if err != nil {
		t.Fatal(err)
	}

	profileSchema := &Schema{
		Name:      "Profile",
		TableName: "profiles",
		Fields: []Field{
			{
				Name:            "ID",
				Type:            valuetype.Int,
				TypeElementExpr: "int",
				Tag: FieldTag{
					Raw:        []string{"primary_key", "immutable"},
					PrimaryKey: true,
					Immutable:  true,
				},
			},
			{
				Name:            "UserID",
				Type:            valuetype.Int,
				TypeElementExpr: "int",
				Tag: FieldTag{
					Raw:   []string{"index"},
					Index: true,
				},
			},
			{
				Name:            "Bio",
				Type:            valuetype.String,
				TypeElementExpr: "string",
				Tag: FieldTag{
					Raw:       []string{"type=text"},
					TableType: "text",
				},
			},
			{
				Name:            "CreatedAt",
				Type:            valuetype.Time,
				TypeElementExpr: "time.Time",
				Tag: FieldTag{
					Raw:       []string{"immutable"},
					Immutable: true,
				},
			},
			{
				Name:            "UpdatedAt",
				Type:            valuetype.Time,
				TypeElementExpr: "time.Time",
				Tag: FieldTag{
					Raw:       []string{"immutable"},
					Immutable: true,
				},
			},
		},
	}

	userSchema := Schema{
		Name:      "User",
		TableName: "users",
		Fields: []Field{
			{
				Name:            "ID",
				Type:            valuetype.Int,
				TypeElementExpr: "int",
				Tag: FieldTag{
					Raw:        []string{"primary_key", "immutable"},
					PrimaryKey: true,
					Immutable:  true,
				},
			},
			{
				Name:            "Username",
				Type:            valuetype.String,
				TypeElementExpr: "string",
				Tag: FieldTag{
					Raw:    []string{"unique"},
					Unique: true,
				},
			},
			{
				Name:            "Email",
				Type:            valuetype.String,
				TypeElementExpr: "string",
				Tag: FieldTag{
					Raw: []string{},
				},
			},
			{
				Name:            "RefreshToken",
				Type:            valuetype.String,
				TypeElementExpr: "string",
				Tag: FieldTag{
					Raw: []string{},
				},
			},
			{
				Name:            "Timezone",
				Type:            valuetype.String,
				TypeElementExpr: "string",
				Tag: FieldTag{
					Raw: []string{},
				},
			},
			{
				Name:            "TimeDiff",
				Type:            valuetype.Int,
				TypeElementExpr: "int",
				Tag: FieldTag{
					Raw: []string{},
				},
			},
			{
				Name:            "CreatedAt",
				Type:            valuetype.Time,
				TypeElementExpr: "time.Time",
				Tag: FieldTag{
					Raw:       []string{"immutable"},
					Immutable: true,
				},
			},
			{
				Name:            "UpdatedAt",
				Type:            valuetype.Time,
				TypeElementExpr: "time.Time",
				Tag: FieldTag{
					Raw:       []string{"immutable"},
					Immutable: true,
				},
			},
		},
	}

	profileField := Field{
		Name:            "Profile",
		Type:            valuetype.Ref(valuetype.FieldValueType("Profile")),
		TypeElementExpr: "Profile",
		Tag: FieldTag{
			Raw:         []string{"association"},
			Association: true,
		},
		Association: &Association{
			Slice:  false,
			Schema: profileSchema,
		},
	}

	postsField := Field{
		Name:            "Posts",
		Type:            valuetype.Slice(valuetype.FieldValueType("Post")),
		TypeElementExpr: "Post",
		Tag: FieldTag{
			Raw:         []string{"association"},
			Association: true,
		},
		Association: &Association{
			Slice: true,
			Schema: &Schema{
				Name:      "Post",
				TableName: "posts",
				Fields: []Field{
					{
						Name:            "ID",
						Type:            valuetype.Int,
						TypeElementExpr: "int",
						Tag: FieldTag{
							Raw:        []string{"primary_key", "immutable"},
							PrimaryKey: true,
							Immutable:  true,
						},
					},
					{
						Name:            "UserID",
						Type:            valuetype.Int,
						TypeElementExpr: "int",
						Tag: FieldTag{
							Raw:   []string{"index"},
							Index: true,
						},
					},
					{
						Name:            "Title",
						Type:            valuetype.String,
						TypeElementExpr: "string",
						Tag: FieldTag{
							Raw: []string{},
						},
					},
					{
						Name:            "Body",
						Type:            valuetype.String,
						TypeElementExpr: "string",
						Tag: FieldTag{
							Raw:       []string{"type=text"},
							TableType: "text",
						},
					},
					{
						Name:            "CreatedAt",
						Type:            valuetype.Time,
						TypeElementExpr: "time.Time",
						Tag: FieldTag{
							Raw:       []string{"immutable"},
							Immutable: true,
						},
					},
					{
						Name:            "UpdatedAt",
						Type:            valuetype.Time,
						TypeElementExpr: "time.Time",
						Tag: FieldTag{
							Raw:       []string{"immutable"},
							Immutable: true,
						},
					},
					{
						Name:            "User",
						Type:            valuetype.FieldValueType("User"),
						TypeElementExpr: "User",
						Tag: FieldTag{
							Raw:         []string{"association"},
							Association: true,
						},
						Association: &Association{
							Slice: false,
							Schema: &Schema{
								Name:      "User",
								TableName: "users",
								Fields: []Field{
									{
										Name:            "ID",
										Type:            valuetype.Int,
										TypeElementExpr: "int",
										Tag: FieldTag{
											Raw:        []string{"primary_key", "immutable"},
											PrimaryKey: true,
											Immutable:  true,
										},
									},
									{
										Name:            "Username",
										Type:            valuetype.String,
										TypeElementExpr: "string",
										Tag: FieldTag{
											Raw:    []string{"unique"},
											Unique: true,
										},
									},
									{
										Name:            "Email",
										Type:            valuetype.String,
										TypeElementExpr: "string",
										Tag: FieldTag{
											Raw: []string{},
										},
									},
									{
										Name:            "RefreshToken",
										Type:            valuetype.String,
										TypeElementExpr: "string",
										Tag: FieldTag{
											Raw: []string{},
										},
									},
									{
										Name:            "Timezone",
										Type:            valuetype.String,
										TypeElementExpr: "string",
										Tag: FieldTag{
											Raw: []string{},
										},
									},
									{
										Name:            "TimeDiff",
										Type:            valuetype.Int,
										TypeElementExpr: "int",
										Tag: FieldTag{
											Raw: []string{},
										},
									},
									{
										Name:            "CreatedAt",
										Type:            valuetype.Time,
										TypeElementExpr: "time.Time",
										Tag: FieldTag{
											Raw:       []string{"immutable"},
											Immutable: true,
										},
									},
									{
										Name:            "UpdatedAt",
										Type:            valuetype.Time,
										TypeElementExpr: "time.Time",
										Tag: FieldTag{
											Raw:       []string{"immutable"},
											Immutable: true,
										},
									},
									{
										Name:            "Profile",
										Type:            valuetype.Ref(valuetype.FieldValueType("Profile")),
										TypeElementExpr: "Profile",
										Tag: FieldTag{
											Raw:         []string{"association"},
											Association: true,
										},
										Association: &Association{
											Slice:  false,
											Schema: profileSchema,
										},
									},
								},
							},
						},
					},
					{
						Name:            "Likes",
						Type:            valuetype.Slice(valuetype.FieldValueType("Like")),
						TypeElementExpr: "Like",
						Tag: FieldTag{
							Raw:         []string{"association"},
							Association: true,
						},
						Association: &Association{
							Slice: true,
							Schema: &Schema{
								Name:      "Like",
								TableName: "likes",
								Fields: []Field{
									{
										Name:            "ID",
										Type:            valuetype.Int,
										TypeElementExpr: "int",
										Tag: FieldTag{
											Raw:        []string{"primary_key", "immutable"},
											PrimaryKey: true,
											Immutable:  true,
										},
									},
									{
										Name:            "LikeableID",
										Type:            valuetype.Int,
										TypeElementExpr: "int",
										Tag: FieldTag{
											Raw:   []string{"index"},
											Index: true,
										},
									},
									{
										Name:            "LikeableType",
										Type:            valuetype.String,
										TypeElementExpr: "string",
										Tag: FieldTag{
											Raw:   []string{"index"},
											Index: true,
										},
									},
									{
										Name:            "CreatedAt",
										Type:            valuetype.Time,
										TypeElementExpr: "time.Time",
										Tag: FieldTag{
											Raw:       []string{"immutable"},
											Immutable: true,
										},
									},
									{
										Name:            "UpdatedAt",
										Type:            valuetype.Time,
										TypeElementExpr: "time.Time",
										Tag: FieldTag{
											Raw:       []string{"immutable"},
											Immutable: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	names := []string{}
	for _, s := range list {
		names = append(names, s.Name)
	}

	if diff := cmp.Diff([]string{"User", "Post", "Profile", "Like"}, names); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	actual := list[0:1]

	profile := actual[0].Fields[8]
	posts := actual[0].Fields[9]
	actual[0].Fields = actual[0].Fields[0:8]
	if diff := cmp.Diff(userSchema, actual[0]); diff != "" {
		t.Errorf("userSchema mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(profileField, profile); diff != "" {
		t.Errorf("profileField mismatch (-want +got):\n%s", diff)
	}

	opt := cmp.FilterValues(func(x, y *Schema) bool {
		return x.Name == "User" || y.Name == "User"
	}, cmp.Ignore())

	if diff := cmp.Diff(postsField, posts, opt); diff != "" {
		t.Errorf("postsField mismatch (-want +got):\n%s", diff)
	}
}
