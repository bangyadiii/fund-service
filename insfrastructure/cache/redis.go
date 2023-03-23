package cache

import (
	"backend-crowdfunding/config"
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type redisClient struct {
	Client *redis.Client
	config config.Config
}

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
}

func Init(config config.Config) RedisClient {
	return &redisClient{
		config: config,
		Client: redis.NewClient(&redis.Options{
			Addr:     config.Get("REDIS_HOST"),
			Password: config.Get("REDIS_PASSWORD"),
			DB:       0,
		}),
	}
}

func (rc *redisClient) Set(ctx context.Context, key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	return rc.Client.WithContext(ctx).Set(key, jsonData, 10*time.Minute).Err()
}

func (rc *redisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.Client.WithContext(ctx).Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, err
}
