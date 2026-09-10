package db

import (
	"database/sql"

	"github.com/ficusinapot/ds/internal/db/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joomcode/errorx"
)

const postgresDriver = "pgx"

type Client struct {
	ent   *ent.Client
	sqlDB *sql.DB
}

func Open(cfg Config) (*Client, error) {
	sqlDB, err := sql.Open(postgresDriver, cfg.DSN)
	if err != nil {
		return nil, ConnectionError.Wrap(err, "open SQL database")
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	driver := entsql.OpenDB(dialect.Postgres, sqlDB)

	return &Client{
		ent:   ent.NewClient(ent.Driver(driver)),
		sqlDB: sqlDB,
	}, nil
}

func (c *Client) Ent() *ent.Client {
	return c.ent
}

func (c *Client) SQLDB() *sql.DB {
	return c.sqlDB
}

func Close(client *Client) error {
	if client == nil || client.ent == nil {
		return nil
	}

	if err := client.ent.Close(); err != nil {
		return errorx.InternalError.Wrap(err, "close sql database")
	}

	return nil
}
