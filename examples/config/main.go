package main

import (
	"os"

	"github.com/4bd1x0/gkit/config"
)

func main() {
	os.Setenv("GO_CONFIG_DIR", "/Users/noname/workspace/4bd1x0/gkit/examples/config/configs")
	config.Init()
	log := config.NewLogger()

	log.Info("Application started")

	ConfigExample()
	LoggerExample()
	DBExample()
	RedisExample()

	log.Info("Application completed")
}
