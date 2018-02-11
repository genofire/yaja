package state

import (
	"crypto/tls"
	"fmt"

	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/model"
	"dev.sum7.eu/genofire/yaja/server/utils"
	"golang.org/x/crypto/acme/autocert"
)

// Start state
type Start struct {
	Next   State
	Client *utils.Client
}

// Process message
func (state *Start) Process() State {
	state.Client.Log = state.Client.Log.WithField("state", "stream")
	state.Client.Log.Debug("running")
	defer state.Client.Log.Debug("leave")

	element, err := state.Client.Read()
	if err != nil {
		state.Client.Log.Warn("unable to read: ", err)
		return nil
	}
	if element.Name.Space != messages.NSStream || element.Name.Local != "stream" {
		state.Client.Log.Warn("is no stream")
		return state
	}
	for _, attr := range element.Attr {
		if attr.Name.Local == "to" {
			state.Client.JID = &model.JID{Domain: attr.Value}
			state.Client.Log = state.Client.Log.WithField("jid", state.Client.JID.Full())
		}
	}
	if state.Client.JID == nil {
		state.Client.Log.Warn("no 'to' domain readed")
		return nil
	}

	fmt.Fprintf(state.Client.Conn, `<?xml version='1.0'?>
		<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>`,
		messages.CreateCookie(), messages.NSClient, messages.NSStream)

	fmt.Fprintf(state.Client.Conn, `<stream:features>
			<starttls xmlns='%s'>
				<required/>
			</starttls>
		</stream:features>`,
		messages.NSStream)

	return state.Next
}

// TLSUpgrade state
type TLSUpgrade struct {
	Next       State
	Client     *utils.Client
	TLSConfig  *tls.Config
	TLSManager *autocert.Manager
}

// Process message
func (state *TLSUpgrade) Process() State {
	state.Client.Log = state.Client.Log.WithField("state", "tls upgrade")
	state.Client.Log.Debug("running")
	defer state.Client.Log.Debug("leave")

	element, err := state.Client.Read()
	if err != nil {
		state.Client.Log.Warn("unable to read: ", err)
		return nil
	}
	if element.Name.Space != messages.NSTLS || element.Name.Local != "starttls" {
		state.Client.Log.Warn("is no starttls", element)
		return nil
	}
	fmt.Fprintf(state.Client.Conn, "<proceed xmlns='%s'/>", messages.NSTLS)
	// perform the TLS handshake
	var tlsConfig *tls.Config
	if m := state.TLSManager; m != nil {
		var cert *tls.Certificate
		cert, err = m.GetCertificate(&tls.ClientHelloInfo{ServerName: state.Client.JID.Domain})
		if err != nil {
			state.Client.Log.Warn("no cert in tls manger found: ", err)
			return nil
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{*cert},
		}
	}
	if tlsConfig == nil {
		tlsConfig = state.TLSConfig
		if tlsConfig != nil {
			tlsConfig.ServerName = state.Client.JID.Domain
		} else {
			state.Client.Log.Warn("no tls config found: ", err)
			return nil
		}
	}

	tlsConn := tls.Server(state.Client.Conn, tlsConfig)
	err = tlsConn.Handshake()
	if err != nil {
		state.Client.Log.Warn("unable to tls handshake: ", err)
		return nil
	}
	// restart the Connection
	state.Client.SetConnecting(tlsConn)

	return state.Next
}
