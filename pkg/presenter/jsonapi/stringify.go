package jsonapi

import (
	"fmt"
	"strconv"
	"time"
)

func Stringify(v any) string {
	if v == nil {
		return "null"
	}

	switch v := v.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case *string:
		if v == nil {
			return "null"
		}

		return fmt.Sprintf("\"%s\"", *v)
	case int:
		return strconv.Itoa(v)
	case *int:
		if v == nil {
			return "null"
		}
		return strconv.Itoa(*v)
	case bool:
		return strconv.FormatBool(v)
	case *bool:
		if v == nil {
			return "null"
		}
		return strconv.FormatBool(*v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case *float64:
		if v == nil {
			return "null"
		}
		return strconv.FormatFloat(*v, 'f', -1, 32)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case *float32:
		if v == nil {
			return "null"
		}
		return strconv.FormatFloat(float64(*v), 'f', -1, 32)
	case time.Time:
		return fmt.Sprintf("\"%s\"", v.Format("2006-01-02T15:04:05 -0700"))
	case *time.Time:
		if v == nil {
			return "null"
		}
		return fmt.Sprintf("\"%s\"", v.Format("2006-01-02T15:04:05 -0700"))
	case fmt.Stringer:
		return fmt.Sprintf("\"%s\"", v)
	default:
		return fmt.Sprintf("\"%s\"", v)
	}
}
