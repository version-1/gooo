package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestMiddleware(t *testing.T) {
	mw := Middlewares{}
	output := []string{}

	mw.Append(Middleware{
		Name: "mw1",
		If:   Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			output = append(output, "mw1")
			return true
		},
	})

	mw.Append(Middleware{
		Name: "mw2",
		If:   Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			output = append(output, "mw2")
			return true
		},
	})

	mw.Append(Middleware{
		Name: "mw3",
		If:   Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			output = append(output, "mw3")
			return true
		},
	})

	mw.Prepend(Middleware{
		Name: "mw5",
		If:   Always,
		Do: func(w http.ResponseWriter, r *http.Request) bool {
			output = append(output, "mw5")
			return true
		},
	})

	expect := []string{"mw5", "mw1", "mw4", "mw2", "mw3"}
	if !reflect.DeepEqual(output, expect) {
		fmt.Printf("order of middlewares is incorrect. expect %v, got %v", expect, output)
	}
}
