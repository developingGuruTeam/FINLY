package redisDB

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis_container:6379",
		Password: "",
		DB:       0,
	})

	// Проверка подключения
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{Client: client}, nil
}
