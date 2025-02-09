package partial

import (
	"fmt"
	"strings"
)

func AnonymousStruct(fields []string) string {
	return fmt.Sprintf(`struct {
    %s
  }`, strings.Join(fields, "\n"))
}
