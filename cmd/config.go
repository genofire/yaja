package cmd

import (
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/model/config"
)

var (
	configPath string
)

func loadConfig() *config.Config {
	config, err := config.ReadConfigFile(configPath)
	if err != nil {
		log.Fatal("unable to load config file:", err)
	}
	return config
}
