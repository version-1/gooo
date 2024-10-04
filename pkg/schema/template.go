package schema

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/version-1/gooo/pkg/errors"
	"github.com/version-1/gooo/pkg/schema/internal/template"
	"github.com/version-1/gooo/pkg/util"
	"golang.org/x/tools/imports"
)

type SchemaTemplate struct {
	filename string
	URL      string
	Package  string
	Schema   Schema
}

func (s SchemaTemplate) Filename() string {
	return fmt.Sprintf("generated--%s", util.Basename(strings.ToLower(s.filename)))
}

func (s SchemaTemplate) Render() (string, error) {
	str := ""
	str += fmt.Sprintf("package %s\n", s.Package)
	str += "\n"

	if len(s.libs()) > 0 {
		str += fmt.Sprintf("import (\n%s\n)\n", strings.Join(s.libs(), "\n"))
	}
	str += "\n"

	// columns
	str += template.Method{
		Receiver:    s.Schema.Name,
		Name:        "Columns",
		Args:        []template.Arg{},
		ReturnTypes: []string{"[]string"},
		Body: fmt.Sprintf(
			"return []string{%s}",
			strings.Join(wrapQuote(s.Schema.Columns()), ", "),
		),
	}.String()

	// scan
	scanFields := []string{}
	for _, f := range s.Schema.ColumnFields() {
		scanFields = append(scanFields, fmt.Sprintf("&obj.%s", f.Name))
	}

	receiver := template.Pointer(s.Schema.Name)
	methods := []template.Method{
		{
			Receiver: receiver,
			Name:     "Scan",
			Args: []template.Arg{
				{Name: "rows", Type: "scanner"},
			},
			ReturnTypes: []string{"error"},
			Body: fmt.Sprintf(`if err := rows.Scan(%s); err != nil {
					return err
				}

				return nil`,
				strings.Join(scanFields, ", "),
			),
		},
		{
			Receiver: receiver,
			Name:     "Destroy",
			Args: []template.Arg{
				{Name: "ctx", Type: "context.Context"},
				{Name: "qr", Type: "queryer"},
			},
			ReturnTypes: []string{"error"},
			Body: fmt.Sprintf(`zero, err := util.IsZero(obj.ID)
			if err != nil {
				return goooerrors.Wrap(err)
			}

      if zero {
			  return goooerrors.Wrap(ErrPrimaryKeyMissing)
			}

			query := "DELETE FROM %s WHERE id = $1"
			if _, err := qr.ExecContext(ctx, query, obj.ID); err != nil {
				return goooerrors.Wrap(err)
			}

			return nil`, s.Schema.TableName),
		},
		{
			Receiver: receiver,
			Name:     "Find",
			Args: []template.Arg{
				{Name: "ctx", Type: "context.Context"},
				{Name: "qr", Type: "queryer"},
			},
			ReturnTypes: []string{"error"},
			Body: fmt.Sprintf(`zero, err := util.IsZero(obj.ID)
			if err != nil {
				return goooerrors.Wrap(err)
			}

			if zero {
			  return goooerrors.Wrap(ErrPrimaryKeyMissing)
			}

			query := "SELECT %s FROM %s WHERE id = $1"
			row := qr.QueryRowContext(ctx, query, obj.ID)

			if err := obj.Scan(row); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return goooerrors.Wrap(ErrNotFound)
				}

				return goooerrors.Wrap(err)
			}

			return nil`,
				strings.Join(s.Schema.Columns(), ", "),
				s.Schema.TableName,
			),
		},
	}

	for _, m := range methods {
		str += m.String()
	}

	str += s.defineSave()
	str += s.defineAssign()
	str += s.defineValidate()
	str += s.defineJSONAPISerialize()
	str += s.defineToJSONAPIResource()

	return pretify(s.Filename(), str)
}

func (s SchemaTemplate) defineValidate() string {
	str := ""
	str += "return nil"

	return template.Method{
		Receiver:    s.Schema.Name,
		Name:        "validate",
		Args:        []template.Arg{},
		ReturnTypes: []string{"ormerrors.ValidationError"},
		Body:        str,
	}.String()
}

func (s SchemaTemplate) defineSave() string {
	query := fmt.Sprintf(`
		INSERT INTO %s (%s) VALUES ($1, $2, $3)
		ON CONFLICT(id) DO UPDATE SET %s
		RETURNING %s
  `,
		s.Schema.TableName,
		strings.Join(s.Schema.MutableColumns(), ", "),
		strings.Join(s.Schema.SetClause(), ", "),
		strings.Join(s.Schema.Columns(), ", "),
	)

	mutableValues := []string{}
	for _, f := range s.Schema.MutableFields() {
		mutableValues = append(mutableValues, fmt.Sprintf("obj.%s", f.Name))
	}

	validateStr := `if err := obj.validate(); err != nil {
				return err
			}
		`

	return template.Method{
		Receiver: template.Pointer(s.Schema.Name),
		Name:     "Save",
		Args: []template.Arg{
			{Name: "ctx", Type: "context.Context"},
			{Name: "qr", Type: "queryer"},
		},
		ReturnTypes: []string{"error"},
		Body: fmt.Sprintf(
			validateStr+
				"query := `%s`\n"+`
			row := qr.QueryRowContext(ctx, query, %s)
			if err := obj.Scan(row); err != nil {
				return err
			}

			return nil`,
			query,
			strings.Join(mutableValues, ", "),
		),
	}.String()
}

func (s SchemaTemplate) defineAssign() string {
	fields := []string{}
	for _, f := range s.Schema.Fields {
		fields = append(fields, fmt.Sprintf("obj.%s = v.%s", f.Name, f.Name))
	}

	return template.Method{
		Receiver: template.Pointer(s.Schema.Name),
		Name:     "Assign",
		Args: []template.Arg{
			{Name: "v", Type: s.Schema.Name},
		},
		ReturnTypes: []string{},
		Body:        strings.Join(fields, "\n"),
	}.String()
}

func (s SchemaTemplate) libs() []string {
	list := []string{
		schemaPackage,
		errorsPackage,
		ormerrPackage,
		stringsPackage,
		jsonapiPackage,
		utilPackage,
		"\"github.com/google/uuid\"",
		"\"strings\"",
		"\"time\"",
		"\"fmt\"",
	}

	return list
}

func wrapQuote(list []string) []string {
	for i := range list {
		list[i] = fmt.Sprintf("\"%s\"", list[i])
	}

	return list
}

func pretify(filename, s string) (string, error) {
	// return s, nil
	formatted, err := format.Source([]byte(s))
	if err != nil {
		return s, errors.Wrap(err)
	}

	processed, err := imports.Process(filename, formatted, nil)
	if err != nil {
		return string(formatted), errors.Wrap(err)
	}

	return string(processed), nil
}
