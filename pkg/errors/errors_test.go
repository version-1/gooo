package errors

import (
	"fmt"
	"strings"
	"testing"

	goootesting "github.com/version-1/gooo/pkg/testing"
)

func TestErrors(t *testing.T) {
	err := New("msg")

	test := goootesting.NewTable([]goootesting.Record[string, []string]{
		{
			Name: "Stacktrace",
			Subject: func(_t *testing.T) string {
				return err.StackTrace()
			},
			Expect: func(t *testing.T) []string {
				return []string{
					"gooo/pkg/errors/errors_test.go. method: TestErrors. line: 12",
					"src/testing/testing.go. method: tRunner. line: 1689",
					"src/runtime/asm_amd64.s. method: goexit. line: 1695",
					"",
				}
			},
			Assert: func(t *testing.T, r *goootesting.Record[string, []string]) bool {
				e := r.Expect(t)
				lines := strings.Split(r.Subject(t), "\n")
				for i, line := range lines {
					if !strings.Contains(line, e[i]) {
						t.Errorf("Expected(line %d) %s to contain %s", i, line, e[i])
						return false
					}
				}
				return true
			},
		},
		{
			Name: "Print Error with +v",
			Subject: func(_t *testing.T) string {
				return fmt.Sprintf("%+v", err)
			},
			Expect: func(t *testing.T) []string {
				return []string{
					"pkg/errors : msg",
					"",
					"gooo/pkg/errors/errors_test.go. method: TestErrors. line: 12",
					"src/testing/testing.go. method: tRunner. line: 1689",
					"src/runtime/asm_amd64.s. method: goexit. line: 1695",
					"",
					"",
				}
			},
			Assert: func(t *testing.T, r *goootesting.Record[string, []string]) bool {
				e := r.Expect(t)
				lines := strings.Split(r.Subject(t), "\n")
				for i, line := range lines {
					if !strings.Contains(line, e[i]) {
						t.Errorf("Expected(line %d) %s to contain %s", i, line, e[i])
						return false
					}
				}
				return true
			},
		},
		{
			Name: "Print Error with v",
			Subject: func(_t *testing.T) string {
				return fmt.Sprintf("%v", err)
			},
			Expect: func(t *testing.T) []string {
				return []string{"pkg/errors : msg"}
			},
			Assert: func(t *testing.T, r *goootesting.Record[string, []string]) bool {
				return r.Subject(t) == r.Expect(t)[0]
			},
		},
		{
			Name: "Print Error with s",
			Subject: func(_t *testing.T) string {
				return fmt.Sprintf("%s", err)
			},
			Expect: func(t *testing.T) []string {
				return []string{"pkg/errors : msg"}
			},
			Assert: func(t *testing.T, r *goootesting.Record[string, []string]) bool {
				return r.Subject(t) == r.Expect(t)[0]
			},
		},
	})

	test.Run(t)
}
