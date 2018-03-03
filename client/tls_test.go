package client

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTLS(t *testing.T) {
	assert := assert.New(t)
	client := &Client{}
	assert.Nil(client.TLSConnectionState())
	client.conn = &tls.Conn{}
	assert.NotNil(client.TLSConnectionState())
}
