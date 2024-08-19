package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

var Client RedisClient

func NewRedisClient(addr, password string, db int) error {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("could not connect to Redis: %v", err)
	}

	Client = RedisClient{
		Client: rdb,
		Ctx:    ctx,
	}

	return nil
}

func (r *RedisClient) Publish(channel, message string) error {
	err := r.Client.Publish(r.Ctx, channel, message).Err()
	if err != nil {
		return fmt.Errorf("could not publish message: %v", err)
	}
	return nil
}
