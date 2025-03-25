package template

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/gooolib/errors"
)

//go:embed components/common/*.go.tmpl
var commonTmpl embed.FS

type CommonPlainTemplateParams struct {
	Package      string
	HeadComments string
	Dependencies []string
	Content      string
}

func (p CommonPlainTemplateParams) Render() (string, error) {
	var b bytes.Buffer
	tmpl := template.Must(template.New("plain").ParseFS(commonTmpl, "components/common/plain.go.tmpl"))
	if err := tmpl.ExecuteTemplate(&b, "plain.go.tmpl", p); err != nil {
		return "", errors.Wrap(err)
	}

	return b.String(), nil
}
