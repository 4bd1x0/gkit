package main

import (
	"github.com/4bd1x0/gkit/config"
)

func main() {
	log := config.NewLogger()

	log.Info("Application started")

	ConfigExample()
	LoggerExample()
	DBExample()
	RedisExample()

	log.Info("Application completed")
}
