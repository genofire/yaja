package client

import (
	"encoding/xml"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

func TestRead(t *testing.T) {
	assert := assert.New(t)

	server, clientConn := net.Pipe()
	client := &Client{}
	client.setConnection(clientConn)

	go server.Write([]byte(`<message>`))

	element, err := client.Read()
	assert.NoError(err)
	assert.Equal("message", element.Name.Local)

	go server.Write([]byte(`<>`))
	element, err = client.Read()
	assert.Error(err)
}

func TestSend(t *testing.T) {
	assert := assert.New(t)

	server, clientConn := net.Pipe()
	client := &Client{
		Logging: log.WithField("test", "send"),
	}
	client.setConnection(clientConn)
	serverDecoder := xml.NewDecoder(server)
	wgWait := &sync.WaitGroup{}

	wgWait.Add(1)
	go func() {
		err := client.Send(3)
		assert.NoError(err)
		wgWait.Done()
	}()

	element, err := read(serverDecoder)
	wgWait.Wait()
	assert.NoError(err)
	assert.Equal("int", element.Name.Local)

	wgWait.Add(1)
	go func() {
		err := client.Send(&xmpp.MessageClient{To: xmppbase.NewJID("a@a.de")})
		assert.NoError(err)
		wgWait.Done()
	}()

	element, err = read(serverDecoder)
	wgWait.Wait()
	assert.NoError(err)
	assert.Equal("message", element.Name.Local)
	assert.Equal("a@a.de", element.Attr[1].Value)

	wgWait.Add(1)
	go func() {
		err := client.Send(&xmpp.IQClient{Type: xmpp.IQTypeGet})
		assert.NoError(err)
		wgWait.Done()
	}()

	element, err = read(serverDecoder)
	wgWait.Wait()
	assert.NoError(err)
	assert.Equal("iq", element.Name.Local)
	assert.Equal("get", element.Attr[2].Value)

	wgWait.Add(1)
	go func() {
		err := client.Send(&xmpp.PresenceClient{Type: xmpp.PresenceTypeSubscribe})
		assert.NoError(err)
		wgWait.Done()
	}()

	element, err = read(serverDecoder)
	wgWait.Wait()
	assert.NoError(err)
	assert.Equal("presence", element.Name.Local)
	assert.Equal("subscribe", element.Attr[1].Value)
}
