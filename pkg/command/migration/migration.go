package migration

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/version-1/gooo/pkg/command/migration/adapter/yaml"
	"github.com/version-1/gooo/pkg/command/migration/constants"
	"github.com/version-1/gooo/pkg/db"
	"github.com/version-1/gooo/pkg/logger"
)

var _ Runner = (*yaml.SchemaManager)(nil)

type Command struct {
	database string
	conn     *db.DB
	runner   Runner
	version  string
	logger   logger.Logger
}

type Runner interface {
	Up(ctx context.Context, db db.Tx, size int) error
	Down(ctx context.Context, db db.Tx, size int) error
}

func NewWith(conn *sqlx.DB, runner Runner, l logger.Logger) (*Command, error) {
	_logger := l
	if _logger == nil {
		_logger = logger.DefaultLogger
	}

	d := db.New(conn)
	c := &Command{
		conn:   d,
		runner: runner,
		logger: _logger,
	}
	if err := c.prepare(); err != nil {
		return nil, err
	}

	database, err := c.Database()
	if err != nil {
		return c, err
	}

	_logger.Infof("connecting database: %s", database)

	return c, nil
}

func (c Command) prepare() error {
	q := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			version VARCHAR(14) NOT NULL PRIMARY KEY,
			cache jsonb,
			kind VARCHAR NOT NULL,
			rollback_query TEXT,
			created_at timestamp NOT NULL default now(),
			updated_at timestamp NOT NULL default now()
		)
	`, constants.ConfigTableName)
	_, err := c.conn.Exec(q)

	return err
}

func (c *Command) Version() (string, error) {
	if c.version != "" {
		return c.version, nil
	}

	q := fmt.Sprintf("SELECT version FROM %s ORDER BY version desc LIMIT 1", constants.ConfigTableName)
	if err := c.conn.QueryRow(q).Scan(&c.version); err != nil {
		return "", err
	}

	return c.version, nil
}

func (c Command) Database() (string, error) {
	if c.database != "" {
		return c.database, nil
	}

	q := "SELECT current_catalog"
	if err := c.conn.QueryRow(q).Scan(&c.database); err != nil {
		return "", err
	}

	return c.database, nil
}

func (c Command) Exec(ctx context.Context, cmd string, args ...string) error {
	getSize := func(args ...string) (int, error) {
		if len(args) > 0 {
			size := args[0]
			return strconv.Atoi(size)
		}
		return 0, nil
	}

	getName := func(args ...string) string {
		if len(args) > 0 {
			return args[0]
		}
		return ""
	}

	switch cmd {
	case "create":
		return c.Create()
	case "drop":
		return c.Drop()
	case "up":
		size, err := getSize(args...)
		if err != nil {
			return err
		}

		return c.Up(ctx, size)
	case "down":
		size, err := getSize(args...)
		if err != nil {
			return err
		}
		return c.Down(ctx, size)
	case "g", "generate":
		name := getName(args...)
		if name == "" {
			return fmt.Errorf("migration name is required")
		}
		return c.Generate(ctx, name)
	default:
		return fmt.Errorf("invalid command: %s", cmd)
	}
}

func (c Command) Create() error {
	c.logger.Infof("Creating database: %s", c.database)
	q := "CREATE DATABASE IF NOT EXISTS " + c.database
	_, err := c.conn.Exec(q)
	return err
}

func (c Command) Drop() error {
	c.logger.Infof("Dropping database: %s", c.database)
	q := "DROP DATABASE IF EXISTS " + c.database
	_, err := c.conn.Exec(q)
	return err
}

func (c Command) Up(ctx context.Context, size int) error {
	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			c.logger.Errorf("recovered error: %+v", r)
			err := tx.Rollback()
			if err != nil {
				c.logger.Errorf("rollback error: %v", err)
			}
		}
	}()

	if size <= 0 {
		c.logger.Infof("Starting migration up")
	} else {
		c.logger.Infof("Starting migration up. size: %d", size)
	}
	if err := c.runner.Up(ctx, tx, size); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c Command) Down(ctx context.Context, size int) error {
	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			c.logger.Errorf("recovered error: %+v", r)
			err := tx.Rollback()
			if err != nil {
				c.logger.Errorf("rollback error: %v", err)
			}
		}
	}()

	if size <= 0 {
		c.logger.Infof("Starting migration down")
	} else {
		c.logger.Infof("Starting migration down. size: %d", size)
	}
	if err := c.runner.Down(context.Background(), tx, size); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c Command) Generate(ctx context.Context, name string) error {
	return fmt.Errorf("not implemented")
}
