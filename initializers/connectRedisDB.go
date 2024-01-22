package initializers

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Client *redis.Client
}

var (
	ErrNil = errors.New("No matching data")
	Ctx    = context.TODO()
)

func ConnectRedisDB(config *Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisDBAddress,
		Password: "",
		DB:       0,
	})

	if err := client.Ping(Ctx).Err(); err != nil {
		return nil, err
	}
	return &Redis{
		Client: client,
	}, nil
}
