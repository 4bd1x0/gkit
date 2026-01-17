package main

import (
	"github.com/4bd1x0/gkit/config"
)

func RedisExample() {
	log := config.NewLogger()

	log.Info("===== Redis Example Start =====")

	cfg := config.GetConfig()

	log.Infof("Redis Host: %s", cfg.Redis.Host)
	log.Infof("Redis Port: %s", cfg.Redis.Port)
	log.Infof("Redis Password: %s", cfg.Redis.Password)
	log.Infof("Redis DB: %d", cfg.Redis.DB)

	redisTimeout := config.GetInt("custom.redis.timeout")
	log.Infof("Redis Timeout from custom config: %d", redisTimeout)

	log.Info("===== Redis Example End =====\n")
}
