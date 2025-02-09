package seeder

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/version-1/gooo/pkg/toolkit/logger"
)

type SeedExecutor struct {
	cfg Config
}

type Logger interface {
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

type Config interface {
	Connstr() string
	Seeders() []Seeder
	Logger() Logger
}

func New(config Config) *SeedExecutor {
	return &SeedExecutor{
		cfg: config,
	}
}

func (s SeedExecutor) logger() logger.Logger {
	return s.cfg.Logger()
}

type Seeder interface {
	Exec(tx *sqlx.Tx) error
}

func (s SeedExecutor) RunWith(tx *sqlx.Tx, name ...string) {
	_name := ""
	if len(name) > 0 {
		_name = name[0]
	}

	for _, seed := range s.cfg.Seeders() {
		seedName := fmt.Sprintf("%T", seed)
		if _name == "" || strings.HasSuffix(seedName, _name) {
			s.logger().Infof("run seed:\t%s\n", seedName)

			err := seed.Exec(tx)
			if err != nil {
				s.logger().Errorf("%s\n", err.Error())
				tx.Rollback()
				panic(err)
			}
		}
	}
	tx.Commit()
}

func (s SeedExecutor) Run(name ...string) {
	db, err := sqlx.Connect("postgres", s.cfg.Connstr())
	if err != nil {
		s.logger().Fatalf(err.Error())
	}
	defer db.Close()

	tx := db.MustBegin()
	s.RunWith(tx, name...)
}
