package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// ReadConfigFile reads a config model from path of a yml file
func ReadConfigFile(path string) (config *Config, err error) {
	config = &Config{}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return
}
