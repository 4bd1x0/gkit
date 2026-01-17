package main

import (
	"github.com/4bd1x0/gkit/config"
)

func ConfigExample() {
	log := config.NewLogger()

	log.Info("===== Config Example Start =====")

	cfg := config.GetConfig()

	log.Infof("DebugMode: %v", cfg.DebugMode)
	log.Infof("Logger Path: %s", cfg.Logger.Path)
	log.Infof("Logger Console: %v", cfg.Logger.Console)
	log.Infof("Logger Level: %s", cfg.Logger.Level)

	log.Infof("Web Host: %s", cfg.Web.Host)
	log.Infof("Web Port: %d", cfg.Web.Port)

	log.Infof("Redis Host: %s", cfg.Redis.Host)
	log.Infof("Redis Port: %s", cfg.Redis.Port)
	log.Infof("Redis DB: %d", cfg.Redis.DB)

	enableNewUI := config.GetBool("custom.feature_flags.enable_new_ui")
	log.Infof("GetBool('custom.feature_flags.enable_new_ui'): %v", enableNewUI)

	enableAnalytics := config.GetBool("custom.feature_flags.enable_analytics")
	log.Infof("GetBool('custom.feature_flags.enable_analytics'): %v", enableAnalytics)

	apiKey := config.GetString("custom.api_keys.third_party")
	log.Infof("GetString('custom.api_keys.third_party'): %s", apiKey)

	customMap := config.GetStringMap("custom")
	log.Infof("custom config map: %+v", customMap)

	log.Info("===== Config Example End =====\n")
}
