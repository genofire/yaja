package config

import "dev.sum7.eu/genofire/yaja/model"

type Config struct {
	TLSDir     string `toml:"tlsdir"`
	PortClient int    `toml:"port_client"`
	PortServer int    `toml:"port_server"`

	Domain []*Domain `toml:"domain"`
}

type Domain struct {
	FQDN       string       `toml:"fqdn"`
	Admins     []*model.JID `toml:"admins"`
	TLSDisable bool         `toml:"tls_disable"`
	TLSPrivate string       `toml:"tls_private"`
	TLSPublic  string       `toml:"tls_public"`
}
