package orm

import "github.com/version-1/gooo/pkg/datasource/orm/errors"

type Model interface {
	Validate() errors.ValidationError
	Fields() []string
	MutableFields() []string
	Values() []any
	TableName() string
	Identifier() string
	Scan(Scanner) (Model, error)
}

type ExtendedModel interface {
	Model
	NewItem() Model
}
