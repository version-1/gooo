package runner

import (
	"fmt"
	"os"
	"path/filepath"
)

func lookUpFileInAncestors(path, filename string) (string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if absPath == "/" {
		return "", fmt.Errorf("file not found in ancestors: %s", filename)
	}

	for _, e := range entries {
		if !e.IsDir() && e.Name() == filename {
			return absPath, nil
		}
	}

	parent := fmt.Sprintf("../%s", path)
	return lookUpFileInAncestors(parent, filename)
}
