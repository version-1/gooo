package schema

import (
	"path/filepath"
	"testing"
)

func TestSchemaCollection_Gen(t *testing.T) {
	dir := "./pkg/schema/internal/schema"

	schemas := SchemaCollection{
		URL:     "github.com/version-1/gooo",
		Package: filepath.Base(dir),
		Dir:     dir,
	}

	if err := schemas.Gen(); err != nil {
		t.Error(err)
	}
}
