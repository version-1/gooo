package migration

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/version-1/gooo/pkg/command/migration/constants"
	"github.com/version-1/gooo/pkg/command/migration/runner"
	"github.com/version-1/gooo/pkg/core/db"
	goooerrors "github.com/version-1/gooo/pkg/toolkit/errors"
	"github.com/version-1/gooo/pkg/toolkit/logger"
)

var _ Runner = (*runner.Yaml)(nil)

type connector interface {
	Connect() (*sqlx.DB, error)
}

type Command struct {
	database  string
	connector connector
	conn      *db.DB
	runner    Runner
	version   string
	logger    logger.Logger
}

type Runner interface {
	Prepare(conn *sqlx.DB) error
	Up(ctx context.Context, db db.Tx, size int) error
	Down(ctx context.Context, db db.Tx, size int) error
	BasePath() string
	Elements() runner.Elements
	Ext() string
}

func NewWith(conn connector, runner Runner, l logger.Logger) (*Command, error) {
	_logger := l
	if _logger == nil {
		_logger = logger.DefaultLogger
	}

	c := &Command{
		connector: conn,
		runner:    runner,
		logger:    _logger,
	}

	return c, nil
}

func (c *Command) connect() error {
	conn, err := c.connector.Connect()
	if err != nil {
		return err
	}

	c.conn = db.New(conn)
	if err := c.prepare(); err != nil {
		conn.Close()
		return err
	}

	if err := c.runner.Prepare(conn); err != nil {
		conn.Close()
		return err
	}

	database, err := c.Database()
	if err != nil {
		conn.Close()
		return err
	}
	c.logger.Infof("connecting database: %s", database)

	c.database = database

	return nil
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

	shouldConnect, err := validateCmd(cmd)
	if err != nil {
		return err
	}

	if shouldConnect {
		if err := c.connect(); err != nil {
			return goooerrors.Wrap(err)
		}
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
	return goooerrors.Wrap(err)
}

func (c Command) Drop() error {
	c.logger.Infof("Dropping database: %s", c.database)
	q := "DROP DATABASE IF EXISTS " + c.database
	_, err := c.conn.Exec(q)
	return goooerrors.Wrap(err)
}

func (c Command) Up(ctx context.Context, size int) error {
	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		return goooerrors.Wrap(err)
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
		return goooerrors.Wrap(err)
	}

	return tx.Commit()
}

func (c Command) Down(ctx context.Context, size int) error {
	tx, err := c.conn.BeginTx(ctx, nil)
	if err != nil {
		return goooerrors.Wrap(err)
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
		return goooerrors.Wrap(err)
	}

	return tx.Commit()
}

func (c Command) Generate(ctx context.Context, name string) error {
	version := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.%s", version, name, c.runner.Ext())
	if name == "initial" {
		filename = fmt.Sprintf("%s_%s.%s", strings.Repeat("0", 14), name, c.runner.Ext())
	}

	path := fmt.Sprintf("%s/%s", c.runner.BasePath(), filename)
	if _, err := os.Stat(path); err == nil {
		return goooerrors.Wrap(fmt.Errorf("migration already exists: %s", path))
	}

	c.logger.Infof("Generating migration path %s", path)
	f, err := os.Create(path)
	if err != nil {
		return goooerrors.Wrap(err)
	}

	defer f.Close()
	return nil
}

func validateCmd(cmd string) (bool, error) {
	candidates := []string{
		"create",
		"drop",
		"up",
		"down",
		"g",
		"generate",
	}

	shouldNotConnect := []string{
		"g",
		"generate",
	}

	for _, c := range candidates {
		if c == cmd {
			for _, s := range shouldNotConnect {
				if s == cmd {
					return false, nil
				}
			}

			return true, nil
		}
	}

	return false, fmt.Errorf("invalid command: %s. expect: [%s]", cmd, strings.Join(candidates, "|"))
}
