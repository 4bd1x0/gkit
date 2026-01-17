package main

import (
	"github.com/4bd1x0/gkit/config"
)

func DBExample() {
	log := config.NewLogger()

	log.Info("===== DB Example Start =====")

	cfg := config.GetConfig()

	if len(cfg.DB) > 0 {
		log.Infof("Found %d database(s):", len(cfg.DB))
		for name, dbCfg := range cfg.DB {
			log.Infof("  DB [%s]:", name)
			log.Infof("    Driver: %s", dbCfg.Driver)
			log.Infof("    Host: %s", dbCfg.Host)
			log.Infof("    Port: %s", dbCfg.Port)
			log.Infof("    User: %s", dbCfg.User)
			log.Infof("    Database: %s", dbCfg.Name)
			log.Infof("    MaxIdle: %d, MaxOpen: %d", dbCfg.MaxIdle, dbCfg.MaxOpen)
		}
	} else {
		log.Info("No database configured")
	}

	if dbCfg, exists := cfg.DB["default"]; exists {
		log.Infof("Default DB: %s://%s:%s/%s",
			dbCfg.Driver, dbCfg.Host, dbCfg.Port, dbCfg.Name)
	}

	if dbCfg, exists := cfg.DB["primary"]; exists {
		log.Infof("Primary DB: %s://%s:%s/%s",
			dbCfg.Driver, dbCfg.Host, dbCfg.Port, dbCfg.Name)
	}

	log.Info("===== DB Example End =====\n")
}
