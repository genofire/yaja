package toserver

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"

	"dev.sum7.eu/genofire/yaja/database"
	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/server/extension"
	"dev.sum7.eu/genofire/yaja/server/state"
	"dev.sum7.eu/genofire/yaja/server/utils"
	"golang.org/x/crypto/acme/autocert"
)

// ConnectionStartup return steps through TCP TLS state
func ConnectionStartup(db *database.State, tlsconfig *tls.Config, tlsmgmt *autocert.Manager, extensions extension.Extensions, c *utils.Client) state.State {
	receiving := &state.ReceivingClient{Extensions: extensions, Client: c}
	sending := &state.SendingClient{Next: receiving, Client: c}
	tlsstream := &TLSStream{
		Next:   sending,
		Client: c,
	}
	tlsupgrade := &state.TLSUpgrade{
		Next:       tlsstream,
		Client:     c,
		TLSConfig:  tlsconfig,
		TLSManager: tlsmgmt,
	}
	dail := &Dailback{
		Next:   tlsupgrade,
		Client: c,
	}
	return &state.Start{Next: dail, Client: c}
}

// TLSStream state
type Dailback struct {
	Next   state.State
	Client *utils.Client
}

// Process messages
func (state *Dailback) Process() state.State {
	state.Client.Log = state.Client.Log.WithField("state", "dialback")
	state.Client.Log.Debug("running")
	defer state.Client.Log.Debug("leave")

	element, err := state.Client.Read()
	if err != nil {
		state.Client.Log.Warn("unable to read: ", err)
		return nil
	}

	// dailback encode
	type dailback struct {
		XMLName xml.Name `xml:"urn:xmpp:ping ping"`
	}
	db := &dailback{}
	if err = state.Client.In.DecodeElement(db, element); err != nil {
		return state.Next
	}

	state.Client.Log.Info(db)
	return state.Next
}

// TLSStream state
type TLSStream struct {
	Next                  state.State
	Client                *utils.Client
	domainRegisterAllowed utils.DomainRegisterAllowed
}

// Process messages
func (state *TLSStream) Process() state.State {
	state.Client.Log = state.Client.Log.WithField("state", "tls stream")
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

	fmt.Fprintf(state.Client.Conn, `<?xml version='1.0'?>
		<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>
		<stream:features>
			<mechanisms xmlns='%s'>
				<mechanism>EXTERNAL</mechanism>
			</mechanisms>
			<bidi xmlns='urn:xmpp:features:bidi'/>
		</stream:features>`,
		utils.CreateCookie(), messages.NSClient, messages.NSStream,
		messages.NSSASL)

	return state.Next
}

// SASLAuth state
type SASLAuth struct {
	Next                  state.State
	Client                *utils.Client
	database              *database.State
	domainRegisterAllowed utils.DomainRegisterAllowed
}

// Process messages
func (state *SASLAuth) Process() state.State {
	state.Client.Log = state.Client.Log.WithField("state", "sasl auth")
	state.Client.Log.Debug("running")
	defer state.Client.Log.Debug("leave")

	// read the full auth stanza
	element, err := state.Client.Read()
	if err != nil {
		state.Client.Log.Warn("unable to read: ", err)
		return nil
	}
	var auth messages.SASLAuth
	if err = state.Client.In.DecodeElement(&auth, element); err != nil {
		return nil
	}
	data, err := base64.StdEncoding.DecodeString(auth.Body)
	if err != nil {
		state.Client.Log.Warn("body decode: ", err)
		return nil
	}

	state.Client.Log.Debug(auth.Mechanism, string(data))

	state.Client.Log.Info("success auth")
	fmt.Fprintf(state.Client.Conn, "<success xmlns='%s'/>", messages.NSSASL)
	return state.Next
}
