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
func ConnectionStartup(db *database.State, tlsconfig *tls.Config, tlsmgmt *autocert.Manager, registerAllowed utils.DomainRegisterAllowed, extensions []extension.Extension) state.State {
	receiving := &ReceivingClient{Extensions: extensions}
	sending := &SendingClient{Next: receiving}
	authedstream := &AuthedStream{Next: sending}
	authedstart := &AuthedStart{Next: authedstream}
	tlsauth := &SASLAuth{
		Next:                  authedstart,
		database:              db,
		domainRegisterAllowed: registerAllowed,
	}
	tlsstream := &TLSStream{
		Next: tlsauth,
		domainRegisterAllowed: registerAllowed,
	}
	return state.ConnectionStartup(tlsstream, tlsconfig, tlsmgmt)
}

// TLSStream state
type TLSStream struct {
	Next                  state.State
	domainRegisterAllowed utils.DomainRegisterAllowed
}

// Process messages
func (state *TLSStream) Process(client *utils.Client) (state.State, *utils.Client) {
	client.Log = client.Log.WithField("state", "tls stream")
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

	fmt.Fprintf(client.Conn, `<?xml version='1.0'?>
		<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>`,
		utils.CreateCookie(), messages.NSClient, messages.NSStream)

	if state.domainRegisterAllowed(client.JID) {
		fmt.Fprintf(client.Conn, `<stream:features>
			<mechanisms xmlns='%s'>
				<mechanism>PLAIN</mechanism>
			</mechanisms>
			<register xmlns='%s'/>
		</stream:features>`,
			messages.NSSASL, messages.NSFeaturesIQRegister)
	} else {
		fmt.Fprintf(client.Conn, `<stream:features>
			<mechanisms xmlns='%s'>
				<mechanism>PLAIN</mechanism>
			</mechanisms>
		</stream:features>`,
			messages.NSSASL)
	}

	return state.Next, client
}

// SASLAuth state
type SASLAuth struct {
	Next                  state.State
	database              *database.State
	domainRegisterAllowed utils.DomainRegisterAllowed
}

// Process messages
func (state *SASLAuth) Process(client *utils.Client) (state.State, *utils.Client) {
	client.Log = client.Log.WithField("state", "sasl auth")
	client.Log.Debug("running")
	defer client.Log.Debug("leave")

	// read the full auth stanza
	element, err := client.Read()
	if err != nil {
		client.Log.Warn("unable to read: ", err)
		return nil, client
	}
	var auth messages.SASLAuth
	if err = client.In.DecodeElement(&auth, element); err != nil {
		client.Log.Info("start substate for registration")
		return &RegisterFormRequest{
			element:               element,
			domainRegisterAllowed: state.domainRegisterAllowed,
			Next: &RegisterRequest{
				domainRegisterAllowed: state.domainRegisterAllowed,
				database:              state.database,
				Next:                  state.Next,
			},
		}, client
	}
	data, err := base64.StdEncoding.DecodeString(auth.Body)
	if err != nil {
		client.Log.Warn("body decode: ", err)
		return nil, client
	}
	info := strings.Split(string(data), "\x00")
	// should check that info[1] starts with client.JID
	client.JID.Local = info[1]
	client.Log = client.Log.WithField("jid", client.JID.Full())
	success, err := state.database.Authenticate(client.JID, info[2])
	if err != nil {
		client.Log.Warn("auth: ", err)
		return nil, client
	}
	if success {
		client.Log.Info("success auth")
		fmt.Fprintf(client.Conn, "<success xmlns='%s'/>", messages.NSSASL)
		return state.Next, client
	}
	client.Log.Warn("failed auth")
	fmt.Fprintf(client.Conn, "<failure xmlns='%s'><not-authorized/></failure>", messages.NSSASL)
	return nil, client

}

// AuthedStart state
type AuthedStart struct {
	Next state.State
}

// Process messages
func (state *AuthedStart) Process(client *utils.Client) (state.State, *utils.Client) {
	client.Log = client.Log.WithField("state", "authed started")
	client.Log.Debug("running")
	defer client.Log.Debug("leave")

	_, err := client.Read()
	if err != nil {
		client.Log.Warn("unable to read: ", err)
		return nil, client
	}
	fmt.Fprintf(client.Conn, `<?xml version='1.0'?>
		<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>`,
		utils.CreateCookie(), messages.NSClient, messages.NSStream)

	fmt.Fprintf(client.Conn, `<stream:features>
			<bind xmlns='%s'/>
		</stream:features>`,
		messages.NSBind)

	return state.Next, client
}

// AuthedStream state
type AuthedStream struct {
	Next state.State
}

// Process messages
func (state *AuthedStream) Process(client *utils.Client) (state.State, *utils.Client) {
	client.Log = client.Log.WithField("state", "authed stream")
	client.Log.Debug("running")
	defer client.Log.Debug("leave")

	// check that it's a bind request
	// read bind request
	element, err := client.Read()
	if err != nil {
		client.Log.Warn("unable to read: ", err)
		return nil, client
	}
	var msg messages.IQ
	if err = client.In.DecodeElement(&msg, element); err != nil {
		client.Log.Warn("is no iq: ", err)
		return nil, client
	}
	if msg.Type != messages.IQTypeSet {
		client.Log.Warn("is no set iq")
		return nil, client
	}
	if msg.Error != nil {
		client.Log.Warn("iq with error: ", msg.Error.Code)
		return nil, client
	}
	type query struct {
		XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
		Resource string   `xml:"resource"`
	}
	q := &query{}
	err = xml.Unmarshal(msg.Body, q)
	if err != nil {
		client.Log.Warn("is no iq bind: ", err)
		return nil, client
	}
	if q.Resource == "" {
		client.JID.Resource = makeResource()
	} else {
		client.JID.Resource = q.Resource
	}
	client.Log = client.Log.WithField("jid", client.JID.Full())
	client.Out.Encode(&messages.IQ{
		Type: messages.IQTypeResult,
		To:   client.JID.String(),
		From: client.JID.Domain,
		ID:   msg.ID,
		Body: []byte(fmt.Sprintf(
			`<bind xmlns='%s'>
				<jid>%s</jid>
			</bind>`,
			messages.NSBind, client.JID.Full())),
	})

	return state.Next, client
}
