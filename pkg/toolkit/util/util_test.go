package util

import (
	"strings"
	"testing"
)

func TestUtilLookupGomodDirPath(t *testing.T) {
	path, err := LookupGomodDirPath()
	if err != nil {
		t.Errorf("Error: %+v", err)
	}

	if !strings.HasSuffix(path, "gooo") {
		t.Errorf("Expected path to end with 'gooo', got %s", path)
	}
}
