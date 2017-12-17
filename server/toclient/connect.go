package toclient

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/genofire/yaja/database"
	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/server/extension"
	"github.com/genofire/yaja/server/state"
	"github.com/genofire/yaja/server/utils"
	"golang.org/x/crypto/acme/autocert"
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
		<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>`,
		utils.CreateCookie(), messages.NSClient, messages.NSStream)

	if state.domainRegisterAllowed(state.Client.JID) {
		fmt.Fprintf(state.Client.Conn, `<stream:features>
			<mechanisms xmlns='%s'>
				<mechanism>PLAIN</mechanism>
			</mechanisms>
			<register xmlns='%s'/>
		</stream:features>`,
			messages.NSSASL, messages.NSFeaturesIQRegister)
	} else {
		fmt.Fprintf(state.Client.Conn, `<stream:features>
			<mechanisms xmlns='%s'>
				<mechanism>PLAIN</mechanism>
			</mechanisms>
		</stream:features>`,
			messages.NSSASL)
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
		fmt.Fprintf(state.Client.Conn, "<success xmlns='%s'/>", messages.NSSASL)
		return state.Next
	}
	state.Client.Log.Warn("failed auth")
	fmt.Fprintf(state.Client.Conn, "<failure xmlns='%s'><not-authorized/></failure>", messages.NSSASL)
	return nil

}

// AuthedStart state
type AuthedStart struct {
	Next   state.State
	Client *utils.Client
}

// Process messages
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
		<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>`,
		utils.CreateCookie(), messages.NSClient, messages.NSStream)

	fmt.Fprintf(state.Client.Conn, `<stream:features>
			<bind xmlns='%s'/>
		</stream:features>`,
		messages.NSBind)

	return state.Next
}

// AuthedStream state
type AuthedStream struct {
	Next   state.State
	Client *utils.Client
}

// Process messages
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
	var msg messages.IQ
	if err = state.Client.In.DecodeElement(&msg, element); err != nil {
		state.Client.Log.Warn("is no iq: ", err)
		return nil
	}
	if msg.Type != messages.IQTypeSet {
		state.Client.Log.Warn("is no set iq")
		return nil
	}
	if msg.Error != nil {
		state.Client.Log.Warn("iq with error: ", msg.Error.Code)
		return nil
	}
	type query struct {
		XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
		Resource string   `xml:"resource"`
	}
	q := &query{}
	err = xml.Unmarshal(msg.Body, q)
	if err != nil {
		state.Client.Log.Warn("is no iq bind: ", err)
		return nil
	}
	if q.Resource == "" {
		state.Client.JID.Resource = makeResource()
	} else {
		state.Client.JID.Resource = q.Resource
	}
	state.Client.Log = state.Client.Log.WithField("jid", state.Client.JID.Full())
	state.Client.Out.Encode(&messages.IQ{
		Type: messages.IQTypeResult,
		To:   state.Client.JID.String(),
		From: state.Client.JID.Domain,
		ID:   msg.ID,
		Body: []byte(fmt.Sprintf(
			`<bind xmlns='%s'>
				<jid>%s</jid>
			</bind>`,
			messages.NSBind, state.Client.JID.Full())),
	})

	return state.Next
}
