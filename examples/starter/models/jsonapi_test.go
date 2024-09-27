package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
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

	uid := []uuid.UUID{
		uuid.MustParse("4018be75-e855-489d-a151-ddb8fc3fd2dc"),
		uuid.MustParse("ccf7a495-ec22-4358-bccd-d77bec8ee037"),
		uuid.MustParse("f7b1b3b4-3b3b-4b3b-8b3b-3b3b3b3b3b3b"),
	}

	postID := []uuid.UUID{
		uuid.MustParse("15fa357d-089d-4816-9924-65a8e2a91eba"),
		uuid.MustParse("e1222719-b9b6-4191-99c6-9b159884f534"),
		uuid.MustParse("17b89f20-d638-4b6a-b732-1b8f08a914d1"),
	}

	users := []User{}
	for i, id := range uid {
		u := NewUser()
		uu := User{
			ID:        id,
			Username:  "test" + strconv.Itoa(i),
			Email:     fmt.Sprintf("test%d@example.com", i),
			CreatedAt: now,
			UpdatedAt: now,
		}
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
					User:      uu,
					Title:     "title" + strconv.Itoa(i),
					Body:      "body" + strconv.Itoa(i),
					Status:    "published",
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

	expectedStr := strings.ReplaceAll(string(expected), "  ", "\t")
	// trailing newline
	expectedStr = string(expectedStr[0 : len(expectedStr)-1])

	if err := diff(expectedStr, s); err != nil {
		fmt.Printf("expect %s\n\n got %s", expectedStr, s)
		t.Fatal(err)
	}
}

func TestResourceSerialize(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2024-08-07T01:58:13+00:00")
	if err != nil {
		t.Fatal(err)
	}

	uid := uuid.MustParse("4018be75-e855-489d-a151-ddb8fc3fd2dc")
	p1 := uuid.MustParse("ccf7a495-ec22-4358-bccd-d77bec8ee037")
	p2 := uuid.MustParse("f7b1b3b4-3b3b-4b3b-8b3b-3b3b3b3b3b3b")
	uu := NewUserWith(User{
		ID:        uid,
		Username:  "test",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	})
	u := NewUserWith(User{
		ID:        uid,
		Username:  "test",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
		Posts: []Post{
			{
				ID:        p1,
				UserID:    uid,
				User:      *uu,
				Title:     "title1",
				Body:      "body1",
				Status:    "published",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        p2,
				UserID:    uid,
				User:      *uu,
				Title:     "title2",
				Body:      "body2",
				Status:    "published",
				CreatedAt: now,
				UpdatedAt: now,
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

	expectedStr := strings.ReplaceAll(string(expected), "  ", "\t")
	// trailing newline
	expectedStr = string(expectedStr[0 : len(expectedStr)-1])

	if err := diff(expectedStr, s); err != nil {
		fmt.Printf("expect %s\n\n got %s", expectedStr, s)
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
