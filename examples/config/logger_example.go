package main

import (
	"github.com/4bd1x0/gkit/config"
	"go.uber.org/zap"
)

func LoggerExample() {
	log := config.NewLogger()

	log.Info("===== Logger Example Start =====")

	// 1. Debug log example (console only, not written to file)
	log.Debug("This is a debug message")
	//log.Debugf("Debug mode is enabled: %v", config.IsDebugging())

	// 2. Info log example (written to info.log + console)
	//cfg := config.GetConfig()
	//log.Infof("Environment: %s", config.GetEnv())
	//log.Infof("Config loaded: DebugMode=%v", cfg.DebugMode)
	//log.Infof("Logger config: Console=%v, Level=%s", cfg.Logger.Console, cfg.Logger.Level)

	// 3. Error log example (written to error.log + console)
	log.Error("This is an error message example")
	log.Errorf("Error example with context: code=%d, msg=%s", 500, "internal error")
	log.Error("Error with fields", zap.String("module", "logger_example"), zap.Int("retry", 3))

	log.Info("===== Logger Example End =====\n")
}
