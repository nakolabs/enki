package infra

import (
	"enuma-elish/config"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Infra struct {
	Postgres *sqlx.DB

	Redis *redis.Client
}

func New(c *config.Config) (*Infra, error) {
	postgres, err := newPostgres(c)
	if err != nil {
		log.Err(err).Msg("postgres error")
		return nil, err
	}

	rdb, err := newRedis(c)
	if err != nil {
		log.Err(err).Msg("redis error")
		return nil, err
	}

	return &Infra{postgres, rdb}, nil
}
