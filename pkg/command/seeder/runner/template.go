package runner

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/version-1/gooo/pkg/command/seeder"
)

var _ seeder.Config = &TemplateRunner{}

type TemplateRunner struct {
	connstr  string
	pathGlob string
	logger   seeder.Logger
	seeders  []seeder.Seeder
}

func NewTemplateRunner(logger seeder.Logger, connstr, pathGlob string) *TemplateRunner {
	paths, err := collectPaths(pathGlob)
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.ParseGlob(pathGlob))
	seeders := []seeder.Seeder{}
	for _, path := range paths {
		seeders = append(seeders, &TemplateSeeder{
			filename: path,
			tmpl:     tmpl,
		})
	}

	return &TemplateRunner{
		logger:   logger,
		connstr:  connstr,
		pathGlob: pathGlob,
		seeders:  seeders,
	}
}

func (t TemplateRunner) Connstr() string {
	return t.connstr
}

func (t TemplateRunner) Logger() seeder.Logger {
	return t.logger
}

func (t TemplateRunner) Seeders() []seeder.Seeder {
	return t.seeders
}

type TemplateSeeder struct {
	filename string
	tmpl     *template.Template
}

func (t TemplateSeeder) Exec(tx *sqlx.Tx) error {
	return t.execWithName(tx, t.filename)
}

func collectPaths(pathGlob string) ([]string, error) {
	rootPath, err := lookUpFileInAncestors(".", "go.mod")
	if err != nil {
		return nil, err
	}

	glob := fmt.Sprintf("%s/%s", rootPath, pathGlob)

	return filepath.Glob(glob)
}

func (t TemplateSeeder) renderTemplate(filename string) (string, error) {
	buf := new(bytes.Buffer)
	if err := t.tmpl.ExecuteTemplate(buf, filename, struct{}{}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (t TemplateSeeder) execWithName(tx *sqlx.Tx, name string) error {
	query, err := t.renderTemplate(name)
	if err != nil {
		return err
	}

	tx.MustExec(query)
	return nil
}
