package state

import (
	"crypto/tls"
	"fmt"

	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/model"
	"github.com/genofire/yaja/server/utils"
	"golang.org/x/crypto/acme/autocert"
)

// ConnectionStartup return steps through TCP TLS state
func ConnectionStartup(after State, tlsconfig *tls.Config, tlsmgmt *autocert.Manager) State {
	tlsupgrade := &TLSUpgrade{
		Next:      after,
		tlsconfig: tlsconfig,
		tlsmgmt:   tlsmgmt,
	}
	stream := &Start{Next: tlsupgrade}
	return stream
}

// Start state
type Start struct {
	Next State
}

// Process message
func (state *Start) Process(client *utils.Client) (State, *utils.Client) {
	client.Log = client.Log.WithField("state", "stream")
	client.Log.Debug("running")
	defer client.Log.Debug("leave")

	element, err := client.Read()
	if err != nil {
		client.Log.Warn("unable to read: ", err)
		return nil, client
	}
	if element.Name.Space != messages.NSStream || element.Name.Local != "stream" {
		client.Log.Warn("is no stream")
		return state, client
	}
	for _, attr := range element.Attr {
		if attr.Name.Local == "to" {
			client.JID = &model.JID{Domain: attr.Value}
			client.Log = client.Log.WithField("jid", client.JID.Full())
		}
	}
	if client.JID == nil {
		client.Log.Warn("no 'to' domain readed")
		return nil, client
	}

	fmt.Fprintf(client.Conn, `<?xml version='1.0'?>
		<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>`,
		utils.CreateCookie(), messages.NSClient, messages.NSStream)

	fmt.Fprintf(client.Conn, `<stream:features>
			<starttls xmlns='%s'>
				<required/>
			</starttls>
		</stream:features>`,
		messages.NSStream)

	return state.Next, client
}

// TLSUpgrade state
type TLSUpgrade struct {
	Next      State
	tlsconfig *tls.Config
	tlsmgmt   *autocert.Manager
}

// Process message
func (state *TLSUpgrade) Process(client *utils.Client) (State, *utils.Client) {
	client.Log = client.Log.WithField("state", "tls upgrade")
	client.Log.Debug("running")
	defer client.Log.Debug("leave")

	element, err := client.Read()
	if err != nil {
		client.Log.Warn("unable to read: ", err)
		return nil, client
	}
	if element.Name.Space != messages.NSTLS || element.Name.Local != "starttls" {
		client.Log.Warn("is no starttls")
		return state, client
	}
	fmt.Fprintf(client.Conn, "<proceed xmlns='%s'/>", messages.NSTLS)
	// perform the TLS handshake
	var tlsConfig *tls.Config
	if m := state.tlsmgmt; m != nil {
		var cert *tls.Certificate
		cert, err = m.GetCertificate(&tls.ClientHelloInfo{ServerName: client.JID.Domain})
		if err != nil {
			client.Log.Warn("no cert in tls manger found: ", err)
			return nil, client
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{*cert},
		}
	}
	if tlsConfig == nil {
		tlsConfig = state.tlsconfig
		if tlsConfig != nil {
			tlsConfig.ServerName = client.JID.Domain
		} else {
			client.Log.Warn("no tls config found: ", err)
			return nil, client
		}
	}

	tlsConn := tls.Server(client.Conn, tlsConfig)
	err = tlsConn.Handshake()
	if err != nil {
		client.Log.Warn("unable to tls handshake: ", err)
		return nil, client
	}
	// restart the Connection
	client.SetConnecting(tlsConn)

	return state.Next, client
}
