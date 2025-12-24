package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context, cfg Config) (*redis.ClusterClient, error) {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{cfg.Address},
		Password: cfg.Password,
	}), nil
}
