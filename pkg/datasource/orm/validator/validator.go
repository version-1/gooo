package validator

import (
	"fmt"
	"regexp"

	"github.com/version-1/gooo/pkg/datasource/orm/errors"
)

type ValidateFunc func(v any) errors.ValidationError

func Required(k string) ValidateFunc {
	return func(v any) errors.ValidationError {
		if v == nil {
			return errors.NewRequiredError(k)
		}

		return nil
	}
}

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

func OneOf(k string, values []fmt.Stringer) ValidateFunc {
	return func(v any) errors.ValidationError {
		for i := range values {
			ele := values[i].String()
			vv := stringify(v)

			if ele != vv {
				return errors.NewMustOneOfError(k, values, vv)
			}
		}

		return nil
	}
}

func Format(k string, r regexp.Regexp) ValidateFunc {
	return func(v any) errors.ValidationError {
		s := stringify(v)
		if r.MatchString(s) {
			return nil
		}

		return errors.NewFormatInvalidError(k, s)
	}
}
