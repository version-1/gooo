package payload

import (
	"bufio"
	"os"
	"strings"
)

type EnvfileLoader[T comparable] struct {
	path    string
	keyMaps map[string]T
}

func NewEnvfileLoader[T comparable](path string) *EnvfileLoader[T] {
	_path := ".env"
	if path != "" {
		_path = path
	}

	return &EnvfileLoader[T]{path: _path}
}

func (l *EnvfileLoader[T]) Load() (*map[string]any, error) {
	f, err := os.Open(l.path)
	if err != nil {
		return nil, err
	}

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	m := &map[string]any{}
	for s.Scan() {
		line := s.Text()
		if len(line) > 0 && line[0] == '#' {
			continue
		}

		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				v := strings.TrimSpace(strings.TrimSuffix(parts[1], "\n"))
				os.Setenv(parts[0], v)
				(*m)[parts[0]] = v
			}
		}
	}

	f.Close()
	return m, nil
}
