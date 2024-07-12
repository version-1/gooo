package errors

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type WrappedError struct {
	err error
}

type withCause interface {
	Cause() error
}

func (e WrappedError) Value() error {
	return e.err
}

func (e WrappedError) Error() string {
	return e.err.Error()
}

func (e WrappedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", e.err)
			return
		}

		if s.Flag('#') {
			if v, ok := e.err.(withCause); ok {
				fmt.Fprintf(s, "%#v", v.Cause())
				return
			}

			fmt.Fprintf(s, "%#v", e.err)
			return
		}

		fmt.Fprintf(s, "%#v", e.err)
	case 's':
		io.WriteString(s, e.err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.err)
	}
}

func Wrap(err error) WrappedError {
	return WrappedError{err: errors.WithStack(err)}
}

func New(msg string) WrappedError {
	return WrappedError{err: errors.WithStack(fmt.Errorf(msg))}
}

func Errorf(format string, args ...interface{}) WrappedError {
	return WrappedError{err: errors.WithStack(fmt.Errorf(format, args...))}
}

func Is(err error, target error) bool {
	if err == nil || target == nil {
		return false
	}

	v, ok := err.(WrappedError)
	if ok {
		return errors.Is(v.Value(), target)
	}
	return errors.Is(err, target)
}

func As(err error, target error) bool {
	return errors.As(err, target)
}
