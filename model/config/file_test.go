package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Function to test the configuration of the microservice
func TestReadConfig(t *testing.T) {
	assert := assert.New(t)

	config, err := ReadConfigFile("../../config_example.conf")
	assert.NoError(err)
	assert.NotNil(config)

	assert.Equal("/tmp", config.TLSDir)

	config, err = ReadConfigFile("../config_example.co")
	assert.Nil(config)
	assert.Contains(err.Error(), "no such file or directory")

	config, err = ReadConfigFile("testdata/config_panic.conf")
	assert.Nil(config)
	assert.Contains(err.Error(), "keys cannot contain")
}
