package toclient

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/acme/autocert"

	"dev.sum7.eu/genofire/yaja/database"
	"dev.sum7.eu/genofire/yaja/server/extension"
	"dev.sum7.eu/genofire/yaja/server/state"
	"dev.sum7.eu/genofire/yaja/server/utils"
	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
	"dev.sum7.eu/genofire/yaja/xmpp/iq"
)

// ConnectionStartup return steps through TCP TLS state
func ConnectionStartup(db *database.State, tlsconfig *tls.Config, tlsmgmt *autocert.Manager, registerAllowed utils.DomainRegisterAllowed, extensions extension.Extensions, c *utils.Client) state.State {
	receiving := &state.ReceivingClient{Extensions: extensions, Client: c}
	sending := &state.SendingClient{Next: receiving, Client: c}
	authedstream := &AuthedStream{Next: sending, Client: c}
	authedstart := &AuthedStart{Next: authedstream, Client: c}
	tlsauth := &SASLAuth{
		Next:                  authedstart,
		Client:                c,
		database:              db,
		domainRegisterAllowed: registerAllowed,
	}
	tlsstream := &TLSStream{
		Next:                  tlsauth,
		Client:                c,
		domainRegisterAllowed: registerAllowed,
	}
	tlsupgrade := &state.TLSUpgrade{
		Next:       tlsstream,
		Client:     c,
		TLSConfig:  tlsconfig,
		TLSManager: tlsmgmt,
	}
	return &state.Start{Next: tlsupgrade, Client: c}
}

// TLSStream state
type TLSStream struct {
	Next                  state.State
	Client                *utils.Client
	domainRegisterAllowed utils.DomainRegisterAllowed
}

// Process xmpp
func (state *TLSStream) Process() state.State {
	state.Client.Log = state.Client.Log.WithField("state", "tls stream")
	state.Client.Log.Debug("running")
	defer state.Client.Log.Debug("leave")

	element, err := state.Client.Read()
	if err != nil {
		state.Client.Log.Warn("unable to read: ", err)
		return nil
	}
	if element.Name.Space != xmpp.NSStream || element.Name.Local != "stream" {
		state.Client.Log.Warn("is no stream")
		return state
	}

	if state.domainRegisterAllowed(state.Client.JID) {
		fmt.Fprintf(state.Client.Conn, `<?xml version='1.0'?>
			<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>
			<stream:features>
				<register xmlns='%s'/>
				<mechanisms xmlns='%s'>
					<mechanism>PLAIN</mechanism>
				</mechanisms>
			</stream:features>`,
			xmpp.CreateCookie(), xmpp.NSClient, xmpp.NSStream,
			xmpp.NSSASL, xmppiq.NSFeatureRegister)
	} else {
		fmt.Fprintf(state.Client.Conn, `<?xml version='1.0'?>
			<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>
			<stream:features>
				<mechanisms xmlns='%s'>
					<mechanism>PLAIN</mechanism>
				</mechanisms>
			</stream:features>`,
			xmpp.CreateCookie(), xmpp.NSClient, xmpp.NSStream,
			xmpp.NSSASL)
	}

	return state.Next
}

// SASLAuth state
type SASLAuth struct {
	Next                  state.State
	Client                *utils.Client
	database              *database.State
	domainRegisterAllowed utils.DomainRegisterAllowed
}

// Process xmpp
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
	var auth xmpp.SASLAuth
	if err = state.Client.In.DecodeElement(&auth, element); err != nil {
		state.Client.Log.Info("start substate for registration")
		return &RegisterFormRequest{
			Next: &RegisterRequest{
				Next:                  state.Next,
				Client:                state.Client,
				database:              state.database,
				domainRegisterAllowed: state.domainRegisterAllowed,
			},
			Client:                state.Client,
			element:               element,
			domainRegisterAllowed: state.domainRegisterAllowed,
		}
	}
	data, err := base64.StdEncoding.DecodeString(auth.Body)
	if err != nil {
		state.Client.Log.Warn("body decode: ", err)
		return nil
	}
	info := strings.Split(string(data), "\x00")
	// should check that info[1] starts with state.Client.JID
	state.Client.JID.Local = info[1]
	state.Client.Log = state.Client.Log.WithField("jid", state.Client.JID.Full())
	success, err := state.database.Authenticate(state.Client.JID, info[2])
	if err != nil {
		state.Client.Log.Warn("auth: ", err)
		return nil
	}
	if success {
		state.Client.Log.Info("success auth")
		fmt.Fprintf(state.Client.Conn, "<success xmlns='%s'/>", xmpp.NSSASL)
		return state.Next
	}
	state.Client.Log.Warn("failed auth")
	fmt.Fprintf(state.Client.Conn, "<failure xmlns='%s'><not-authorized/></failure>", xmpp.NSSASL)
	return nil

}

// AuthedStart state
type AuthedStart struct {
	Next   state.State
	Client *utils.Client
}

// Process xmpp
func (state *AuthedStart) Process() state.State {
	state.Client.Log = state.Client.Log.WithField("state", "authed started")
	state.Client.Log.Debug("running")
	defer state.Client.Log.Debug("leave")

	_, err := state.Client.Read()
	if err != nil {
		state.Client.Log.Warn("unable to read: ", err)
		return nil
	}
	fmt.Fprintf(state.Client.Conn, `<?xml version='1.0'?>
		<stream:stream xmlns:stream='%s' xml:lang='en' from='%s' id='%x' version='1.0' xmlns='%s'>
		<stream:features>
				<bind xmlns='%s'>
					<required/>
				</bind>
			</stream:features>`,
		xmpp.NSStream, state.Client.JID.Domain, xmpp.CreateCookie(), xmpp.NSClient,
		xmpp.NSBind)

	return state.Next
}

// AuthedStream state
type AuthedStream struct {
	Next   state.State
	Client *utils.Client
}

// Process xmpp
func (state *AuthedStream) Process() state.State {
	state.Client.Log = state.Client.Log.WithField("state", "authed stream")
	state.Client.Log.Debug("running")
	defer state.Client.Log.Debug("leave")

	// check that it's a bind request
	// read bind request
	element, err := state.Client.Read()
	if err != nil {
		state.Client.Log.Warn("unable to read: ", err)
		return nil
	}
	var msg xmpp.IQClient
	if err = state.Client.In.DecodeElement(&msg, element); err != nil {
		state.Client.Log.Warn("is no iq: ", err)
		return nil
	}
	if msg.Type != xmpp.IQTypeSet {
		state.Client.Log.Warn("is no set iq")
		return nil
	}
	if msg.Error != nil {
		state.Client.Log.Warn("iq with error: ", msg.Error.Code)
		return nil
	}

	if msg.Bind == nil {
		state.Client.Log.Warn("is no iq bind: ", err)
		return nil
	}
	if msg.Bind.Resource == "" {
		state.Client.JID.Resource = makeResource()
	} else {
		state.Client.JID.Resource = msg.Bind.Resource
	}
	state.Client.Log = state.Client.Log.WithField("jid", state.Client.JID.Full())
	state.Client.Out.Encode(&xmpp.IQClient{
		Type: xmpp.IQTypeResult,
		To:   state.Client.JID,
		From: xmppbase.NewJID(state.Client.JID.Domain),
		ID:   msg.ID,
		Bind: &xmpp.Bind{JID: state.Client.JID},
	})

	return state.Next
}
