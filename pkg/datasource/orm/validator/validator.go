package validator

import (
	"fmt"
	"regexp"

	"github.com/version-1/gooo/pkg/datasource/orm/errors"
)

type Validator func(k string) ValidatorFunc
type ValidatorFunc func(v ...any) errors.ValidationError

func Required(k string) ValidatorFunc {
	return func(v ...any) errors.ValidationError {
		if v == nil {
			return errors.NewRequiredError(k)
		}

		return nil
	}
}

func OneOf(values []fmt.Stringer) Validator {
	return func(key string) ValidatorFunc {
		return func(v ...any) errors.ValidationError {
			for i := range values {
				ele := values[i].String()
				vv := stringify(v)

				if ele != vv {
					return errors.NewMustOneOfError(key, values, vv)
				}
			}

			return nil
		}
	}
}

func Format(r *regexp.Regexp) Validator {
	return func(k string) ValidatorFunc {
		return func(v ...any) errors.ValidationError {
			if len(v) == 0 {
				return nil
			}

			s := stringify(v[0])
			if r.MatchString(s) {
				return nil
			}

			return errors.NewFormatInvalidError(k, s)
		}
	}
}

var Email = Format(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`))

func stringify(v any) string {
	switch vv := v.(type) {
	case string:
		return vv
	case int:
		return fmt.Sprintf("%d", vv)
	case fmt.Stringer:
		return vv.String()
	default:
		return fmt.Sprintf("%s", vv)
	}
}
