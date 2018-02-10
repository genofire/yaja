package tester

import (
	"dev.sum7.eu/genofire/yaja/model"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	TLSDir       string       `toml:"tlsdir"`
	AccountsPath string       `toml:"accounts_path"`
	OutputPath   string       `toml:"output_path"`
	Logging      log.Level    `toml:"logging"`
	Webserver    string       `toml:"webserver"`
	Admins       []*model.JID `toml:"admins"`
	Client       struct {
		JID      *model.JID `toml:"jid"`
		Password string     `toml:"password"`
	} `toml:"client"`
}
