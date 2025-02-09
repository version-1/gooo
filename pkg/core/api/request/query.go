package request

import (
	"net/url"
	"strconv"

	"github.com/version-1/gooo/pkg/toolkit/logger"
)

type Query struct {
	url    url.URL
	logger logger.Logger
}

func (q Query) GetString(key string) (string, bool) {
	v := q.url.Query().Get(key)
	return v, v != ""
}

func (q Query) GetInt(key string) (int, bool) {
	v := q.url.Query().Get(key)
	if v == "" {
		return 0, false
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		q.logger.Errorf("failed to convert query param %s to int: %s", key, err)
		return 0, false
	}

	return i, true
}
