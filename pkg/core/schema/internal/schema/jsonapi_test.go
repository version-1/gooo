package fixtures

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/version-1/gooo/pkg/presenter/jsonapi"
)

type Meta struct {
	Total   int
	Page    int
	HasNext bool
	HasPrev bool
}

func (m Meta) JSONAPISerialize() (string, error) {
	data := map[string]any{
		"total":    m.Total,
		"page":     m.Page,
		"has_next": m.HasNext,
		"has_prev": m.HasPrev,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func TestResourcesSerialize(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2024-08-07T01:58:13+00:00")
	if err != nil {
		t.Fatal(err)
	}

	uid := []int{
		1,
		2,
		3,
	}

	postID := []int{
		4,
		5,
		6,
	}

	users := []User{}
	for i, id := range uid {
		u := NewUser()
		u.Assign(User{
			ID:        id,
			Username:  "test" + strconv.Itoa(i),
			Email:     fmt.Sprintf("test%d@example.com", i),
			CreatedAt: now,
			UpdatedAt: now,
			Posts: []Post{
				{
					ID:        postID[i],
					UserID:    id,
					Title:     "title" + strconv.Itoa(i),
					Body:      "body" + strconv.Itoa(i),
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		})

		users = append(users, *u)
	}
	root, err := jsonapi.NewManyFrom(
		users,
		Meta{
			Total:   3,
			Page:    1,
			HasNext: true,
			HasPrev: true,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	s, err := root.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile("./fixtures/test_resources_serialize.json")
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	if err := json.Compact(buf, expected); err != nil {
		t.Fatal(err)
	}

	if err := diff(buf.String(), s); err != nil {
		fmt.Printf("expect %s\n\n got %s \n\n\n", buf.String(), s)
		t.Fatal(err)
	}
}

func TestResourceSerialize(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2024-08-07T01:58:13+00:00")
	if err != nil {
		t.Fatal(err)
	}

	uid := 1
	p1 := 10
	p2 := 11
	u := NewUserWith(User{
		ID:           uid,
		Username:     "test",
		Email:        "test@example.com",
		RefreshToken: "refresh_token",
		Timezone:     "Asia/Tokyo",
		TimeDiff:     9,
		CreatedAt:    now,
		UpdatedAt:    now,
		Posts: []Post{
			{
				ID:        p1,
				UserID:    uid,
				Title:     "title1",
				Body:      "body1",
				CreatedAt: now,
				UpdatedAt: now,
				User: User{
					ID:           uid,
					Username:     "test",
					Email:        "test@example.com",
					RefreshToken: "refresh_token",
					Timezone:     "Asia/Tokyo",
					TimeDiff:     9,
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
			{
				ID:        p2,
				UserID:    uid,
				Title:     "title2",
				Body:      "body2",
				CreatedAt: now,
				UpdatedAt: now,
				User: User{
					ID:           uid,
					Username:     "test",
					Email:        "test@example.com",
					RefreshToken: "refresh_token",
					Timezone:     "Asia/Tokyo",
					TimeDiff:     9,
					CreatedAt:    now,
					UpdatedAt:    now,
				},
			},
		},
	})

	resource, includes := u.ToJSONAPIResource()

	root, err := jsonapi.New(resource, includes, nil)
	if err != nil {
		t.Fatal(err)
	}

	s, err := root.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile("./fixtures/test_resource_serialize.json")
	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}
	if err := json.Compact(buf, expected); err != nil {
		t.Fatal(err)
	}

	if err := diff(buf.String(), s); err != nil {
		fmt.Printf("expect %s\n\n got %s\n\n", buf.String(), s)
		t.Fatal(err)
	}
}

func diff(expected, got string) error {
	line := 1
	for i := 0; i < len(expected); i++ {
		if i >= len(got) {
			return errors.New(fmt.Sprintf("got diff at %d line %d. expected(%d), but got(%d)", i, line, len(expected), len(got)))
		}

		if expected[i] != got[i] {
			expectedLines := strings.Split(expected, "\n")
			gotLines := strings.Split(got, "\n")
			msg := fmt.Sprintf("got diff at %d line %d. expected \"%s\", but got \"%s\"", i, line, string(expected[i]), string(got[i]))
			if line > 1 {
				msg += fmt.Sprintf("  %s\n", expectedLines[line-1-1])
			}
			msg += fmt.Sprintf("- %s\n", expectedLines[line-1])
			if line < len(expectedLines) {
				msg += fmt.Sprintf("- %s\n", expectedLines[line])
			}
			msg += "\n\n\n"
			msg += fmt.Sprintf("+ %s\n", gotLines[line-1])
			return errors.New(msg)
		}

		if expected[i] == '\n' {
			line++
		}
	}

	return nil
}
