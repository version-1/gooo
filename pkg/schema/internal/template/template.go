package template

import (
	"fmt"
	"strings"
)

func Pointer(name string) string {
	return "*" + name
}

type Method struct {
	Receiver    string
	Name        string
	Args        []Arg
	ReturnTypes []string
	Body        string
}

func (m Method) String() string {
	return fmt.Sprintf(
		"func (obj %s) %s(%s) (%s) {\n%s\n}\n\n",
		m.Receiver,
		m.Name,
		stringifyArgs(m.Args),
		strings.Join(m.ReturnTypes, ", "),
		m.Body,
	)
}

func stringifyArgs(args []Arg) string {
	str := []string{}
	for _, a := range args {
		str = append(str, a.String())
	}

	return strings.Join(str, ", ")
}

type Arg struct {
	Name string
	Type string
}

func (a Arg) String() string {
	return fmt.Sprintf("%s %s", a.Name, a.Type)
}

type Interface struct {
	Name   string
	Inters []string
}

func (i Interface) String() string {
	str := fmt.Sprintf("type %s interface {\n", i.Name)
	for _, i := range i.Inters {
		str += fmt.Sprintf("\t%s\n", i)
	}
	str += "}\n"

	return str
}
