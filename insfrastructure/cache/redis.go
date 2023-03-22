package cache

import (
	"backend-crowdfunding/config"
	"github.com/go-redis/redis"
)

type Redis struct {
	Client *redis.Client
	config config.Config
}

func Init(config config.Config) Redis {
	return Redis{
		config: config,
		Client: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		}),
	}
}
