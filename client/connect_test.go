package client

import (
	"encoding/xml"
	"fmt"
	"net"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

func TestStartStream(t *testing.T) {
	assert := assert.New(t)

	server, clientConn := net.Pipe()
	client := &Client{
		JID:     xmppbase.NewJID("a@example.com"),
		Logging: log.WithField("test", "startStream"),
	}
	client.setConnection(clientConn)
	wgWait := &sync.WaitGroup{}

	// complete connection
	wgWait.Add(1)
	go func() {
		decoder := xml.NewDecoder(server)
		elm, inlineErr := read(decoder)
		assert.NoError(inlineErr)
		assert.Equal("http://etherx.jabber.org/streams", elm.Name.Space)

		_, inlineErr = fmt.Fprintf(server, "<?xml version='1.0'?>\n"+
			"<stream:stream to='%s' xmlns='%s'\n"+
			" xmlns:stream='%s' version='1.0'>\n",
			"example.com", xmpp.NSClient, xmpp.NSStream)
		assert.NoError(inlineErr)

		_, inlineErr = server.Write([]byte(`
			<features xmlns="http://etherx.jabber.org/streams">
				<starttls xmlns="urn:ietf:params:xml:ns:xmpp-tls"/>
				<bind xmlns="urn:ietf:params:xml:ns:xmpp-bind"/>
				<mechanisms xmlns="urn:ietf:params:xml:ns:xmpp-sasl">
					<mechanism>PLAIN</mechanism>
					<mechanism>notworking</mechanism>
				</mechanisms>
			</features>
			`))
		assert.NoError(inlineErr)
		wgWait.Done()
	}()

	_, err := client.startStream()
	wgWait.Wait()
	assert.NoError(err)

	// no features first rechieve
	wgWait.Add(1)
	go func() {
		decoder := xml.NewDecoder(server)
		elm, inlineErr := read(decoder)
		assert.NoError(inlineErr)
		assert.Equal("http://etherx.jabber.org/streams", elm.Name.Space)

		_, inlineErr = fmt.Fprintf(server, "<?xml version='1.0'?>\n"+
			"<stream:stream to='%s' xmlns='%s'\n"+
			" xmlns:stream='%s' version='1.0'>\n",
			"example.com", xmpp.NSClient, xmpp.NSStream)
		assert.NoError(inlineErr)

		_, inlineErr = server.Write([]byte(`
			<f>
			`))
		assert.NoError(inlineErr)
		wgWait.Done()
	}()

	_, err = client.startStream()
	wgWait.Wait()
	assert.Error(err)
	assert.Contains(err.Error(), "<features>")

	// no stream receive
	wgWait.Add(1)
	go func() {
		decoder := xml.NewDecoder(server)
		elm, inlineErr := read(decoder)
		assert.NoError(inlineErr)
		assert.Equal("http://etherx.jabber.org/streams", elm.Name.Space)

		_, inlineErr = fmt.Fprintf(server, "<s>")
		assert.NoError(inlineErr)

		wgWait.Done()
	}()

	_, err = client.startStream()
	wgWait.Wait()
	assert.Error(err)
	assert.Contains(err.Error(), "is no stream")

	// client  disconnect after stream start
	wgWait.Add(1)
	go func() {
		decoder := xml.NewDecoder(server)
		elm, inlineErr := read(decoder)
		assert.NoError(inlineErr)
		assert.Equal("http://etherx.jabber.org/streams", elm.Name.Space)

		client.Close()

		wgWait.Done()
	}()

	_, err = client.startStream()
	wgWait.Wait()
	assert.Error(err)
	assert.Contains(err.Error(), "closed pipe")

	// client  disconnect before stream start
	_, err = client.startStream()
	wgWait.Wait()
	assert.Error(err)
	assert.Contains(err.Error(), "closed pipe")
}
