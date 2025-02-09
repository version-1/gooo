package errors

import (
	"database/sql"
	"errors"
	"fmt"
)

type DefaultValidationError struct {
	key string
	err error
}

func (d DefaultValidationError) Key() string {
	return d.key
}

func (d DefaultValidationError) Error() string {
	return d.err.Error()
}

var PointerModelExpectedErr = PointerModelExpectedError{}

func NewPointerModelExpectedError(m any) PointerModelExpectedError {
	return PointerModelExpectedError{typeName: fmt.Sprintf("%T", m)}
}

type PointerModelExpectedError struct {
	typeName string
}

func (p PointerModelExpectedError) Error() string {
	return fmt.Sprintf("pointer model expected %s", p.typeName)
}

type ValidationError interface {
	Key() string
	Error() string
}

var _ ValidationError = &DefaultValidationError{}
var _ ValidationError = &NotStructError{}
var _ ValidationError = &FormatInvalidError{}
var _ ValidationError = &RequiredError{}
var _ ValidationError = &MustOneOfError{}

func NewValidationError(key string, v string) ValidationError {
	return DefaultValidationError{key, errors.New(v)}
}

func NewFormatInvalidError(key string, v string) *FormatInvalidError {
	return &FormatInvalidError{key: key, v: v}
}

func NewRequiredError(key string) *RequiredError {
	return &RequiredError{key: key}
}

func NewMustOneOfError(key string, values []fmt.Stringer, value string) *MustOneOfError {
	return &MustOneOfError{key: key, values: values, value: value}
}

func NewNotStructError(v any) *NotStructError {
	return &NotStructError{v: v}
}

type NotStructError struct {
	key string
	v   any
}

func (n NotStructError) Key() string {
	return n.key
}

func (n NotStructError) Error() string {
	return fmt.Sprintf("not struct: %T", n.v)
}

type FormatInvalidError struct {
	key string
	v   string
}

func (f FormatInvalidError) Error() string {
	return fmt.Sprintf("format invalid: %s", f.v)
}

func (f FormatInvalidError) Key() string {
	return f.key
}

type RequiredError struct {
	key string
}

func (f RequiredError) Key() string {
	return f.key
}

func (f RequiredError) Error() string {
	return fmt.Sprintf("can not be empty")
}

type MustOneOfError struct {
	key    string
	values []fmt.Stringer
	value  string
}

func (f MustOneOfError) Key() string {
	return f.key
}

func (f MustOneOfError) Error() string {
	return fmt.Sprintf("must one of %s but got %s", f.values, f.value)
}

type NotFoundError struct {
	err error
}

func (e NotFoundError) Error() string {
	return e.err.Error()
}

var ErrNotFound = NotFoundError{sql.ErrNoRows}
