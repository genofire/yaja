package config

type Config struct {
	TLSDir     string `toml:"tlsdir"`
	StatePath  string `toml:"state_path"`
	PortClient int    `toml:"port_client"`
	PortServer int    `toml:"port_server"`
}
