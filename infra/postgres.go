package infra

import (
	"enuma-elish/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

func newPostgres(c *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.Postgres.Username, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.Database, c.Postgres.SSLMode)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.Postgres.MaxOpenConn)
	db.SetMaxIdleConns(c.Postgres.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(c.Postgres.MaxConnLifetime) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(c.Postgres.MaxConnIdleTime) * time.Minute)

	return db, nil
}
