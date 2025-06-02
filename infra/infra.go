package infra

import (
	"enuma-elish/config"
	"enuma-elish/pkg/cloudinary"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Infra struct {
	Postgres   *sqlx.DB
	Redis      *redis.Client
	Cloudinary *cloudinary.Service
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

	// Initialize Cloudinary
	cloudinaryService, err := cloudinary.New(
		c.Cloudinary.CloudName,
		c.Cloudinary.APIKey,
		c.Cloudinary.APISecret,
		c.Cloudinary.Folder,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &Infra{
		Postgres:   postgres,
		Redis:      rdb,
		Cloudinary: cloudinaryService,
	}, nil
}
