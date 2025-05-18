package infra

import (
	"context"
	"enuma-elish/config"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func newRedis(c *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port),
		Username: c.Redis.Username,
		Password: c.Redis.Password,
		DB:       c.Redis.Database,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
