package seeders

import (
	"github.com/jmoiron/sqlx"
	"github.com/version-1/gooo/pkg/command/seeder"
	"github.com/version-1/gooo/pkg/logger"
)

type DevelopmentSeed struct {
	connstr string
}

func NewDevelopmentSeed(connstr string) DevelopmentSeed {
	return DevelopmentSeed{
		connstr: connstr,
	}
}

func (s DevelopmentSeed) Connstr() string {
	return s.connstr
}

func (s DevelopmentSeed) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		Seed_0001_User{},
	}
}

func (S DevelopmentSeed) Logger() seeder.Logger {
	return logger.DefaultLogger
}

type Seed_0001_User struct{}

func (s Seed_0001_User) Exec(tx *sqlx.Tx) error {
	query := "INSERT INTO seeder_users (name, email) VALUES ('John Doe', 'john@example.com')"
	if _, err := tx.Exec(query); err != nil {
		return err
	}

	return nil
}
