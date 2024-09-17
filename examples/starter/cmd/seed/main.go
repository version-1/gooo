package seed

import (
	"os"

	"github.com/version-1/gooo/examples/starter/db/seeders"
	"github.com/version-1/gooo/pkg/command/seeder"
)

func main() {
	seed := seeders.NewDevelopmentSeed(os.Getenv("DATABASE_URL"))

	ex := seeder.New(seed)
	ex.Run()
}
