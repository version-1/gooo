package orm

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/version-1/gooo/pkg/generator"
	"golang.org/x/tools/imports"
)

type SchemaCollection struct {
	Dir     string
	Package string
	Schemas []Schema
}

func (s SchemaCollection) Gen() error {
	g := generator.Generator{
		Dir:      s.Dir,
		Template: s,
	}

	if err := g.Run(); err != nil {
		return err
	}

	for _, schema := range s.Schemas {
		tmpl := SchemaTemplate{
			filename: schema.Name,
			Package:  s.Package,
			Schema:   schema,
		}

		g := generator.Generator{
			Dir:      s.Dir,
			Template: tmpl,
		}

		if err := g.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (s SchemaCollection) Filename() string {
	return "shared"
}

func (s SchemaCollection) Render() (string, error) {
	str := ""
	str += fmt.Sprintf("package %s\n", s.Package)
	str += "\n"
	str += "// this file is generated by gooo ORM. DON'T EDIT this file\n"

	sharedLibs := wrapQuote([]string{
		"context",
		"database/sql",
		"time",
	})

	if len(sharedLibs) > 0 {
		str += fmt.Sprintf("import (\n%s\n)\n", strings.Join(sharedLibs, "\n"))
	}
	str += "\n"

	str += defineInterface("scanner", []string{
		"Scan(dest ...any) error",
	})

	str += defineInterface("queryer", []string{
		"QueryRowContext(ctx context.Context, query string, dest ...any) *sql.Row",
		"QueryContext(ctx context.Context, query string, dest ...any) (*sql.Rows, error)",
		"ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)",
	})

	return pretify(s.Filename(), str)
}

type SchemaTemplate struct {
	filename string
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

	if len(libs()) > 0 {
		str += fmt.Sprintf("import (\n%s\n)\n", strings.Join(libs(), "\n"))
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
			  return errors.New("primaryKey is required")
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
			  return errors.New("primaryKey is required")
			}

			query := "SELECT %s FROM %s WHERE id = $1"
			row := qr.QueryRowContext(ctx, query, obj.ID)

			if err := obj.Scan(row); err != nil {
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

	return pretify(s.Filename(), str)
}

type Method struct {
	Pointer     bool
	Name        string
	Args        []Arg
	ReturnTypes []string
	Body        string
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

	return s.defineMethod(
		true,
		"Save", []Arg{
			{"ctx", "context.Context"},
			{"qr", "queryer"},
		},
		[]string{"error"},
		fmt.Sprintf("query := `%s`\n"+`
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

func (s SchemaTemplate) defineFind() string {
	return s.defineMethod(
		true,
		"Find", []Arg{
			{"ctx", "context.Context"},
			{"qr", "queryer"},
		},
		[]string{"error"},
		fmt.Sprintf(`if obj.ID == uuid.Nil {
			  return errors.New("primaryKey is required")
			}

			query := "SELECT %s FROM %s WHERE id = $1"
			row := qr.QueryRowContext(ctx, query, obj.ID)

			if err := obj.Scan(row); err != nil {
				return err
			}

			return nil`,
			strings.Join(s.Schema.Columns(), ", "),
			s.Schema.TableName,
		),
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

func libs() []string {
	list := []string{
		"context",
		"time",
		"github.com/google/uuid",
	}

	return wrapQuote(list)
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
		return s, err
	}

	processed, err := imports.Process(filename, formatted, nil)
	if err != nil {
		return string(formatted), err
	}

	return string(processed), err
}
