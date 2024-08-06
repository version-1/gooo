package schema

import (
	"fmt"
	"go/format"
	"strings"

	"golang.org/x/tools/imports"
)

type SchemaTemplate struct {
	filename string
	URL      string
	Package  string
	Schema   Schema
}

func (s SchemaTemplate) Filename() string {
	return strings.ToLower(s.filename)
}

func (s SchemaTemplate) Render() (string, error) {
	str := ""
	str += fmt.Sprintf("package %s\n", s.Package)
	str += "\n"

	if len(s.libs()) > 0 {
		str += fmt.Sprintf("import (\n%s\n)\n", strings.Join(s.libs(), "\n"))
	}
	str += "\n"

	// define type
	str += s.defineStruct()

	// columns
	str += s.defineMethod(
		false,
		"Columns",
		[]Arg{},
		[]string{"[]string"},
		fmt.Sprintf(
			"return []string{%s}",
			strings.Join(wrapQuote(s.Schema.Columns()), ", "),
		),
	)

	// scan
	scanFields := []string{}
	for _, f := range s.Schema.Fields {
		scanFields = append(scanFields, fmt.Sprintf("&obj.%s", f.Name))
	}

	methods := []Method{
		{
			Pointer: true,
			Name:    "Scan",
			Args: []Arg{
				{"rows", "scanner"},
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
			Pointer: true,
			Name:    "Destroy",
			Args: []Arg{
				{"ctx", "context.Context"},
				{"qr", "queryer"},
			},
			ReturnTypes: []string{"error"},
			Body: fmt.Sprintf(`if obj.ID == uuid.Nil {
			  return ErrPrimaryKeyMissing
			}

			query := "DELETE FROM %s WHERE id = $1"
			if _, err := qr.ExecContext(ctx, query, obj.ID); err != nil {
				return err
			}

			return nil`, s.Schema.TableName),
		},
		{
			Pointer: true,
			Name:    "Find",
			Args: []Arg{
				{"ctx", "context.Context"},
				{"qr", "queryer"},
			},
			ReturnTypes: []string{"error"},
			Body: fmt.Sprintf(`if obj.ID == uuid.Nil {
			  return ErrPrimaryKeyMissing
			}

			query := "SELECT %s FROM %s WHERE id = $1"
			row := qr.QueryRowContext(ctx, query, obj.ID)

			if err := obj.Scan(row); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return ErrNotFound
				}

				return err
			}

			return nil`,
				strings.Join(s.Schema.Columns(), ", "),
				s.Schema.TableName,
			),
		},
	}

	for _, m := range methods {
		str += s.defineMethod(m.Pointer, m.Name, m.Args, m.ReturnTypes, m.Body)
	}

	str += s.defineSave()
	str += s.defineAssign()
	str += s.defineValidate()

	return pretify(s.Filename(), str)
}

type Method struct {
	Pointer     bool
	Name        string
	Args        []Arg
	ReturnTypes []string
	Body        string
}

func (s SchemaTemplate) defineValidate() string {
	fields := s.Schema.Fields
	str := ""
	index := 0
	for i, f := range fields {
		for j, validator := range f.Options.Validators {
			values := []string{
				"obj." + f.Name,
			}
			for _, v := range validator.Fields {
				values = append(values, fmt.Sprintf("obj.%s", v))
			}

			if index == 0 {
				str += fmt.Sprintf(`validator := obj.Schema.Fields[%d].Options.Validators[%d]`+"\n", i, j)
			} else {
				str += fmt.Sprintf(`validator = obj.Schema.Fields[%d].Options.Validators[%d]`+"\n", i, j)
			}
			str += fmt.Sprintf(`if err := validator.Validate("%s")(%s); err != nil {
					return err
				}
				`+"\n\n", f.Name, strings.Join(values, ", "))
			index++
		}
	}

	str += "return nil"

	return s.defineMethod(
		false,
		"validate",
		[]Arg{},
		[]string{"goooerrors.ValidationError"},
		str,
	)
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

	return s.defineMethod(
		true,
		"Save", []Arg{
			{"ctx", "context.Context"},
			{"qr", "queryer"},
		},
		[]string{"error"},
		fmt.Sprintf(
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
	)
}

func (s SchemaTemplate) defineAssign() string {
	str := ""
	for _, f := range s.Schema.Fields {
		str += fmt.Sprintf("obj.%s = v.%s\n", f.Name, f.Name)
	}

	return s.defineMethod(
		true,
		"Assign", []Arg{
			{"v", s.Schema.Name},
		},
		[]string{},
		str,
	)
}

func defineInterface(name string, inters []string) string {
	str := fmt.Sprintf("type %s interface {\n", name)
	for _, i := range inters {
		str += fmt.Sprintf("\t%s\n", i)
	}
	str += "}\n"

	return str
}

func (s SchemaTemplate) defineStruct() string {
	str := fmt.Sprintf("type %s struct {\n", s.Schema.Name)
	str += "schema.Schema\n"
	for _, f := range s.Schema.Fields {
		field := fmt.Sprintf("\t%s %s", f.Name, f.Type)
		if f.Tag != "" {
			str += fmt.Sprintf("%s `%s`\n", field, f.Tag)
		} else {
			str += fmt.Sprintf("%s\n", field)
		}
	}
	str += "}\n"
	str += "\n"

	return str
}

type Arg struct {
	Name string
	Type string
}

func (a Arg) String() string {
	return fmt.Sprintf("%s %s", a.Name, a.Type)
}

func stringifyArgs(args []Arg) string {
	str := []string{}
	for _, a := range args {
		str = append(str, a.String())
	}

	return strings.Join(str, ", ")
}

func (s SchemaTemplate) defineMethod(pointer bool, name string, args []Arg, returnTypes []string, body string) string {
	receiver := s.Schema.Name
	if pointer {
		receiver = "*" + receiver
	}

	str := fmt.Sprintf("func (obj %s) %s (%s) (%s) {\n", receiver, name, stringifyArgs(args), strings.Join(returnTypes, ", "))
	str += body
	str += "\n"
	str += "}\n"
	str += "\n"

	return str
}

func (s SchemaTemplate) libs() []string {
	list := []string{
		schemaPackage,
		errorsPackage,
		"\"github.com/google/uuid\"",
		// fmt.Sprintf("schema \"%s/schema\"", s.URL),
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
	formatted, err := format.Source([]byte(s))
	if err != nil {
		return s, err
	}

	processed, err := imports.Process(filename, formatted, nil)
	if err != nil {
		return string(formatted), err
	}

	return string(processed), err
}
