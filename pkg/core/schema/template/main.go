package template

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/gooolib/errors"
	"github.com/version-1/gooo/pkg/core/schema/openapi/v3_0_0"
)

//go:embed components/*.go.tmpl
var tmpl embed.FS

type Main struct {
	Schema       *v3_0_0.RootSchema
	Dependencies []string
}

func (m Main) Filename() string {
	return "cmd/main"
}

func (m Main) Render() (string, error) {
	tmpl := template.Must(template.New("entry").ParseFS(tmpl, "components/entry.go.tmpl"))
	var b bytes.Buffer
	if err := tmpl.ExecuteTemplate(&b, "entry.go.tmpl", m); err != nil {
		return "", err
	}

	res, err := pretify(m.Filename(), b.String())
	if err != nil {
		return "", errors.Wrap(err)
	}
	return string(res), err
}
