package tester

import (
	"dev.sum7.eu/genofire/yaja/messages"
	"github.com/FreifunkBremen/yanic/lib/duration"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	TLSDir         string            `toml:"tlsdir"`
	AccountsPath   string            `toml:"accounts_path"`
	OutputPath     string            `toml:"output_path"`
	Logging        log.Level         `toml:"logging"`
	LoggingClients log.Level         `toml:"logging_clients"`
	LoggingBots    log.Level         `toml:"logging_bots"`
	Timeout        duration.Duration `toml:"timeout"`
	Interval       duration.Duration `toml:"interval"`
	Admins         []*messages.JID   `toml:"admins"`
	Client         struct {
		JID      *messages.JID `toml:"jid"`
		Password string        `toml:"password"`
	} `toml:"client"`
}
