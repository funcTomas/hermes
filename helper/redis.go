package helper

import (
	"context"
	"log"

	"github.com/funcTomas/hermes/common/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	cfg *config.RedisConfig
}

func NewRedisClient(cfg *config.RedisConfig) *RedisClient {
	return &RedisClient{cfg: cfg}
}

func (c *RedisClient) Connect(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.cfg.Addr,
		Password: c.cfg.Password,
		DB:       c.cfg.DB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *RedisClient) Close(client *redis.Client) error {
	log.Println("Closing Redis connection")
	return client.Close()
}
