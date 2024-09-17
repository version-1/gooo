package db

import (
	"fmt"
	"strings"

	"github.com/version-1/gooo/pkg/logger"
)

type QueryLogger interface {
	Log(query string, args ...any)
	Println(v ...any)
	Printf(format string, v ...any)
}

type queryLoggerAdapter struct {
	logger logger.Logger
}

func (q *queryLoggerAdapter) Log(query string, args ...any) {
	_args := []string{}
	_query := query
	for i, arg := range args {
		_args = append(_args, humanize(arg))
		_query = strings.Replace(_query, fmt.Sprintf("$%d", i+1), humanize(arg), 1)
	}

	q.logger.Infof("executing query: %s args: %s", _query, strings.Join(_args, ", "))
}

func (q *queryLoggerAdapter) Println(v ...any) {
	q.logger.Infof("%s", v...)
}

func (q *queryLoggerAdapter) Printf(format string, v ...any) {
	q.logger.Infof(format, v...)
}

var defaultLogger = &queryLoggerAdapter{logger: logger.DefaultLogger}

func humanize(v any) string {
	switch vv := v.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case int:
		return fmt.Sprintf("%s", v)
	case []int:
		list := []string{}
		for _, item := range vv {
			list = append(list, fmt.Sprintf("%d", item))
		}

		return fmt.Sprintf("[%s]", strings.Join(list, ", "))
	case []string:
		list := []string{}
		for _, item := range vv {
			list = append(list, fmt.Sprintf("'%s'", item))
		}

		return fmt.Sprintf("[%s]", strings.Join(list, ", "))
	case *string:
		if vv == nil {
			return "NULL"
		}
		return fmt.Sprintf("'%s'", v)
	case *int:
		if vv == nil {
			return "NULL"
		}
		return fmt.Sprintf("%s", v)
	case *[]int:
		if vv == nil {
			return "NULL"
		}

		list := []string{}
		for _, item := range *vv {
			list = append(list, fmt.Sprintf("%d", item))
		}

		return fmt.Sprintf("[%s]", strings.Join(list, ", "))
	case *[]string:
		if vv == nil {
			return "NULL"
		}

		list := []string{}
		for _, item := range *vv {
			list = append(list, fmt.Sprintf("'%s'", item))
		}

		return fmt.Sprintf("[%s]", strings.Join(list, ", "))
	case nil:
		return "NULL"
	case fmt.GoStringer:
		return fmt.Sprintf("'%s'", vv.GoString())
	case fmt.Stringer:
		return fmt.Sprintf("'%s'", vv.String())
	default:
		return fmt.Sprintf("%v", v)
	}

	return ""
}
