package config

type Config struct {
	TLSDir    string `toml:"tlsdir"`
	StatePath string `toml:"state_path"`
	Address   struct {
		Client []string `toml:"client"`
		Server []string `toml:"server"`
	} `toml:"address"`
}
