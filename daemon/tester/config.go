package tester

import (
	"dev.sum7.eu/genofire/yaja/model"
	"github.com/FreifunkBremen/yanic/lib/duration"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	TLSDir       string            `toml:"tlsdir"`
	AccountsPath string            `toml:"accounts_path"`
	OutputPath   string            `toml:"output_path"`
	Logging      log.Level         `toml:"logging"`
	Timeout      duration.Duration `toml:"timeout"`
	Interval     duration.Duration `toml:"interval"`
	Admins       []*model.JID      `toml:"admins"`
	Client       struct {
		JID      *model.JID `toml:"jid"`
		Password string     `toml:"password"`
	} `toml:"client"`
}
