package logging

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type driver interface {
	Infof(string, ...any)
}

type QueryLogger struct {
	driver driver
}

func (l QueryLogger) Info(query string, args ...any) {
	s := renderQuery(query, args...)
	l.driver.Infof("Query: %s\n Args: %s\n", s, args)
}

func NewQueryLogger(driver driver) *QueryLogger {
	return &QueryLogger{driver: driver}
}

func renderQuery(query string, args ...any) string {
	q := query
	for i, a := range args {
		tmpl, value := resolveAny(a)

		q = strings.ReplaceAll(q, fmt.Sprintf("$%d", i+1), fmt.Sprintf(tmpl, value))
	}

	return q
}

const defaultTruncate = 100

func stringifySlice(slice any) string {
	res := []string{}
	g := reflect.ValueOf(slice)

	for i := 0; i < g.Len(); i++ {
		s := g.Index(i).Interface()
		switch v := s.(type) {
		case fmt.Stringer:
			res = append(res, fmt.Sprintf("'%s'", truncate(v.String(), defaultTruncate)))
		case int:
			res = append(res, strconv.Itoa(v))
		default:
			res = append(res, fmt.Sprintf("'%s'", truncate(v, defaultTruncate)))
		}
	}

	return "ARRAY[" + strings.Join(res, ", ") + "]"
}

func resolveAny(a any) (string, any) {
	tmpl := "%s"
	var value any
	switch v := a.(type) {
	case bool:
		tmpl = "%t"
		value = v
	case int:
		tmpl = "%d"
		value = v
	case []string:
		tmpl = "%s"
		value = stringifySlice(a)
	case fmt.Stringer, string:
		tmpl = "'%s'"
		value = truncate(v, defaultTruncate)
	default:
		tmpl = "%s"
		value = truncate(v, defaultTruncate)
	}

	return tmpl, value
}

func truncate(s any, n int) string {
	v := fmt.Sprintf("%s", s)
	if len(v) > n {
		return v[:n] + "..."
	}

	return v
}
