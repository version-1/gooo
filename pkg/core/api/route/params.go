package route

import (
	"fmt"
	"strconv"
	"strings"
)

type Params struct {
	m map[string]string
}

func parseParams(matcher, path string) Params {
	p := Params{m: make(map[string]string)}

	pathSegments := strings.Split(path, "/")
	matcherSegments := strings.Split(matcher, "/")
	for i, part := range matcherSegments {
		if strings.HasPrefix(part, ":") {
			if len(pathSegments) > i {
				p.m[part] = pathSegments[i]
			}
		}
	}

	return p
}

func (p Params) GetBool(key string) (bool, error) {
	v, err := p.GetString(key)
	if err != nil {
		return false, err
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, err
	}

	return b, nil
}

func (p Params) GetString(key string) (string, error) {
	v, ok := p.m[key]
	if !ok {
		return "", fmt.Errorf("param %s not found", key)
	}

	return v, nil
}

func (p Params) GetInt(key string) (int, error) {
	v, err := p.GetString(key)
	if err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return n, nil
}
