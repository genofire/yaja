package server

import (
	log "github.com/sirupsen/logrus"
)

type Config struct {
	TLSDir    string `toml:"tlsdir"`
	StatePath string `toml:"state_path"`
	Logging   struct {
		Level       log.Level `toml:"level"`
		LevelClient log.Level `toml:"level_client"`
		LevelServer log.Level `toml:"level_server"`
	} `toml:"logging"`
	Register struct {
		Enable  bool     `toml:"enable"`
		Domains []string `toml:"domains"`
	} `toml:"register"`
	Address struct {
		Webserver []string `toml:"webserver"`
		Client    []string `toml:"client"`
		Server    []string `toml:"server"`
	} `toml:"address"`
}
