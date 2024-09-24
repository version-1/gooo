package errors

import (
	"fmt"
	"strings"
	"testing"
)

func TestErrors(t *testing.T) {
	err := New("msg")

	tests := []struct {
		name    string
		subject func() string
		expect  string
	}{
		{
			name: "Stacktrace",
			subject: func() string {
				return err.StackTrace()
			},
			expect: strings.Join([]string{
				"/Users/admin/Projects/Private/gooo/pkg/errors/errors_test.go:github.com/version-1/gooo/pkg/errors.TestErrors:10",
				"/usr/local/Cellar/go/1.22.3/libexec/src/testing/testing.go:testing.tRunner:1689",
				"/usr/local/Cellar/go/1.22.3/libexec/src/runtime/asm_amd64.s:runtime.goexit:1695",
			}, "\n") + "\n",
		},
		{
			name: "Print Error with +v",
			subject: func() string {
				return fmt.Sprintf("%+v", err)
			},
			expect: strings.Join([]string{
				"/Users/admin/Projects/Private/gooo/pkg/errors/errors_test.go:github.com/version-1/gooo/pkg/errors.TestErrors:10",
				"/usr/local/Cellar/go/1.22.3/libexec/src/testing/testing.go:testing.tRunner:1689",
				"/usr/local/Cellar/go/1.22.3/libexec/src/runtime/asm_amd64.s:runtime.goexit:1695",
			}, "\n") + "\n",
		},
		{
			name: "Print Error with v",
			subject: func() string {
				return fmt.Sprintf("%v", err)
			},
			expect: "pkg/errors : msg",
		},
		{
			name: "Print Error with s",
			subject: func() string {
				return fmt.Sprintf("%s", err)
			},
			expect: "pkg/errors : msg",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.subject(); strings.TrimSpace(got) != strings.TrimSpace(test.expect) {
				t.Errorf("expected\n %s, got\n %s", test.expect, got)
				fmt.Printf("expected len %d, expected got %d", len(test.expect), len(got))

				for i, c := range test.expect {
					if string(c) != string(got[i]) {
						t.Errorf("%d. expected %s, got %s", i, string(c), string(got[i]))
					}
				}
			}
		})
	}
}
