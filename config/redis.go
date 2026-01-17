package config

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

func NewRedisClient() *redis.Client {
	redisOnce.Do(func() {
		cfg := &config.Redis
		host := cfg.Host
		if host == "" {
			host = "127.0.0.1"
		}
		port := cfg.Port
		if port == "" {
			port = "6379"
		}
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: cfg.Password,
			DB:       cfg.DB,
		})
	})
	return redisClient
}

// CloseRedisClient closes Redis connection (call on program exit)
func CloseRedisClient() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}
