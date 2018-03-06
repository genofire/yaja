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

func TestDecode(t *testing.T) {
	assert := assert.New(t)

	server, clientConn := net.Pipe()
	client := &Client{
		Logging: log.WithField("test", "decode"),
	}
	client.setConnection(clientConn)

	go server.Write([]byte(`<message xmlns="jabber:client" to="a@example.com"></message>`))

	msg := &xmpp.MessageClient{}
	err := client.ReadDecode(msg)
	assert.NoError(err)
	assert.Equal("a@example.com", msg.To.String())

	go server.Write([]byte(`<iq xmlns="jabber:client"  to="a@example.com"></iq>`))

	iq := &xmpp.IQClient{}
	err = client.ReadDecode(iq)
	assert.NoError(err)
	assert.Equal("a@example.com", iq.To.String())
	assert.Nil(iq.Ping)

	go server.Write([]byte(`<iq xmlns="jabber:client" type="result"><ping xmlns="urn:xmpp:ping"/></iq>`))

	err = client.ReadDecode(iq)
	assert.NoError(err)
	assert.NotNil(iq.Ping)

	wgWait := &sync.WaitGroup{}
	go server.Write([]byte(`<iq xmlns="jabber:client" type="get"><ping xmlns="urn:xmpp:ping"/></iq>`))
	wgWait.Add(1)
	go func() {
		_, inlineErr := read(xml.NewDecoder(server))
		assert.NoError(inlineErr)
		wgWait.Done()
	}()

	err = client.ReadDecode(iq)
	wgWait.Wait()
	assert.NoError(err)

	go server.Write([]byte(`<>`))

	err = client.ReadDecode(msg)
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
		inlineErr := client.Send(&xmpp.MessageClient{To: xmppbase.NewJID("a@a.de")})
		assert.NoError(inlineErr)
		wgWait.Done()
	}()

	element, err = read(serverDecoder)
	wgWait.Wait()
	assert.NoError(err)
	assert.Equal("message", element.Name.Local)
	assert.Equal("a@a.de", element.Attr[1].Value)

	wgWait.Add(1)
	go func() {
		inlineErr := client.Send(&xmpp.IQClient{Type: xmpp.IQTypeGet})
		assert.NoError(inlineErr)
		wgWait.Done()
	}()

	element, err = read(serverDecoder)
	wgWait.Wait()
	assert.NoError(err)
	assert.Equal("iq", element.Name.Local)
	assert.Equal("get", element.Attr[2].Value)

	wgWait.Add(1)
	go func() {
		inlineErr := client.Send(&xmpp.PresenceClient{Type: xmpp.PresenceTypeSubscribe})
		assert.NoError(inlineErr)
		wgWait.Done()
	}()

	element, err = read(serverDecoder)
	wgWait.Wait()
	assert.NoError(err)
	assert.Equal("presence", element.Name.Local)
	assert.Equal("subscribe", element.Attr[1].Value)
}
