package schemav2

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed components/*.go.tmpl
var tmpl embed.FS

type Main struct {
	Routes string
}

func (m Main) Filename() string {
	return "main"
}

func (m Main) Render() (string, error) {
	tmpl := template.Must(template.New("entry").ParseFS(tmpl, "components/entry.go.tmpl"))
	var b bytes.Buffer
	if err := tmpl.ExecuteTemplate(&b, "entry.go.tmpl", m); err != nil {
		return "", err
	}

	res, err := pretify(m.Filename(), b.String())
	return string(res), err
}
