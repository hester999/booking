package redis

import (
	"booking/internal/config"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg config.CacheConfig) (*RedisClient, error) {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:6379/%s", cfg.Address, cfg.Database))
	if err != nil {
		return nil, err
	}

	if opt == nil {
		return nil, errors.New("redis options is nil")
	}

	client := redis.NewClient(opt)

	return &RedisClient{
		Client: client,
	}, nil
}

func (r *RedisClient) Close() error {
	if r == nil || r.Client == nil {
		return nil
	}
	return r.Client.Close()
}
