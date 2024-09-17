package seed

import (
	"os"

	"github.com/version-1/gooo/pkg/command/seeder"
	"github.com/version-1/gooo/pkg/command/seeder/runner"
	"github.com/version-1/gooo/pkg/logger"
)

func main() {
	tmpl := runner.NewTemplateRunner(logger.DefaultLogger, os.Getenv("DATABAE_URL"), "db/seeders/template/*.sql")

	ex := seeder.New(tmpl)
	ex.Run()
}
