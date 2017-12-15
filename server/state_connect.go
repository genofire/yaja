package server

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/model"
)

// ConnectionStartup return steps through TCP TLS state
func ConnectionStartup() State {
	receiving := &ReceivingClient{}
	sending := &SendingClient{Next: receiving}
	authedstream := &AuthedStream{Next: sending}
	authedstart := &AuthedStart{Next: authedstream}
	tlsauth := &SASLAuth{Next: authedstart}
	tlsstream := &TLSStream{Next: tlsauth}
	tlsupgrade := &TLSUpgrade{Next: tlsstream}
	stream := &Start{Next: tlsupgrade}
	return stream
}

// Start state
type Start struct {
	Next State
}

// Process message
func (state *Start) Process(client *Client) (State, *Client) {
	client.log = client.log.WithField("state", "stream")
	client.log.Debug("running")
	defer client.log.Debug("leave")

	element, err := client.Read()
	if err != nil {
		client.log.Warn("unable to read: ", err)
		return nil, client
	}
	if element.Name.Space != messages.NSStream || element.Name.Local != "stream" {
		client.log.Warn("is no stream")
		return state, client
	}
	for _, attr := range element.Attr {
		if attr.Name.Local == "to" {
			client.jid = &model.JID{Domain: attr.Value}
			client.log = client.log.WithField("jid", client.jid.Full())
		}
	}
	if client.jid == nil {
		client.log.Warn("no 'to' domain readed")
		return nil, client
	}

	fmt.Fprintf(client.Conn, `<?xml version='1.0'?>
		<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>`,
		createCookie(), messages.NSClient, messages.NSStream)

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
	Next State
}

// Process message
func (state *TLSUpgrade) Process(client *Client) (State, *Client) {
	client.log = client.log.WithField("state", "tls upgrade")
	client.log.Debug("running")
	defer client.log.Debug("leave")

	element, err := client.Read()
	if err != nil {
		client.log.Warn("unable to read: ", err)
		return nil, client
	}
	if element.Name.Space != messages.NSTLS || element.Name.Local != "starttls" {
		client.log.Warn("is no starttls")
		return state, client
	}
	fmt.Fprintf(client.Conn, "<proceed xmlns='%s'/>", messages.NSTLS)
	// perform the TLS handshake
	var tlsConfig *tls.Config
	if m := client.Server.TLSManager; m != nil {
		var cert *tls.Certificate
		cert, err = m.GetCertificate(&tls.ClientHelloInfo{ServerName: client.jid.Domain})
		if err != nil {
			client.log.Warn("no cert in tls manger found: ", err)
			return nil, client
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{*cert},
		}
	}
	if tlsConfig == nil {
		tlsConfig = client.Server.TLSConfig
		if tlsConfig != nil {
			tlsConfig.ServerName = client.jid.Domain
		} else {
			client.log.Warn("no tls config found: ", err)
			return nil, client
		}
	}

	tlsConn := tls.Server(client.Conn, tlsConfig)
	err = tlsConn.Handshake()
	if err != nil {
		client.log.Warn("unable to tls handshake: ", err)
		return nil, client
	}
	// restart the Connection
	client.NewConnecting(tlsConn)

	return state.Next, client
}

// TLSStream state
type TLSStream struct {
	Next State
}

// Process messages
func (state *TLSStream) Process(client *Client) (State, *Client) {
	client.log = client.log.WithField("state", "tls stream")
	client.log.Debug("running")
	defer client.log.Debug("leave")

	element, err := client.Read()
	if err != nil {
		client.log.Warn("unable to read: ", err)
		return nil, client
	}
	if element.Name.Space != messages.NSStream || element.Name.Local != "stream" {
		client.log.Warn("is no stream")
		return state, client
	}

	fmt.Fprintf(client.Conn, `<?xml version='1.0'?>
		<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>`,
		createCookie(), messages.NSClient, messages.NSStream)

	if client.DomainRegisterAllowed() {
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
	Next State
}

// Process messages
func (state *SASLAuth) Process(client *Client) (State, *Client) {
	client.log = client.log.WithField("state", "sasl auth")
	client.log.Debug("running")
	defer client.log.Debug("leave")

	// read the full auth stanza
	element, err := client.Read()
	if err != nil {
		client.log.Warn("unable to read: ", err)
		return nil, client
	}
	var auth messages.SASLAuth
	if err = client.in.DecodeElement(&auth, element); err != nil {
		client.log.Info("start substate for registration")
		return &RegisterFormRequest{
			element: element,
			Next: &RegisterRequest{
				Next: state.Next,
			},
		}, client
	}
	data, err := base64.StdEncoding.DecodeString(auth.Body)
	if err != nil {
		client.log.Warn("body decode: ", err)
		return nil, client
	}
	info := strings.Split(string(data), "\x00")
	// should check that info[1] starts with client.jid
	client.jid.Local = info[1]
	client.log = client.log.WithField("jid", client.jid.Full())
	success, err := client.Server.Database.Authenticate(client.jid, info[2])
	if err != nil {
		client.log.Warn("auth: ", err)
		return nil, client
	}
	if success {
		client.log.Info("success auth")
		fmt.Fprintf(client.Conn, "<success xmlns='%s'/>", messages.NSSASL)
		return state.Next, client
	}
	client.log.Warn("failed auth")
	fmt.Fprintf(client.Conn, "<failure xmlns='%s'><not-authorized/></failure>", messages.NSSASL)
	return nil, client

}

// AuthedStart state
type AuthedStart struct {
	Next State
}

// Process messages
func (state *AuthedStart) Process(client *Client) (State, *Client) {
	client.log = client.log.WithField("state", "authed started")
	client.log.Debug("running")
	defer client.log.Debug("leave")

	_, err := client.Read()
	if err != nil {
		client.log.Warn("unable to read: ", err)
		return nil, client
	}
	fmt.Fprintf(client.Conn, `<?xml version='1.0'?>
		<stream:stream id='%x' version='1.0' xmlns='%s' xmlns:stream='%s'>`,
		createCookie(), messages.NSClient, messages.NSStream)

	fmt.Fprintf(client.Conn, `<stream:features>
			<bind xmlns='%s'/>
		</stream:features>`,
		messages.NSBind)

	return state.Next, client
}

// AuthedStream state
type AuthedStream struct {
	Next State
}

// Process messages
func (state *AuthedStream) Process(client *Client) (State, *Client) {
	client.log = client.log.WithField("state", "authed stream")
	client.log.Debug("running")
	defer client.log.Debug("leave")

	// check that it's a bind request
	// read bind request
	element, err := client.Read()
	if err != nil {
		client.log.Warn("unable to read: ", err)
		return nil, client
	}
	var msg messages.IQ
	if err = client.in.DecodeElement(&msg, element); err != nil {
		client.log.Warn("is no iq: ", err)
		return nil, client
	}
	if msg.Type != messages.IQTypeSet {
		client.log.Warn("is no set iq")
		return nil, client
	}
	if msg.Error != nil {
		client.log.Warn("iq with error: ", msg.Error.Code)
		return nil, client
	}
	type query struct {
		XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
		Resource string   `xml:"resource"`
	}
	q := &query{}
	err = xml.Unmarshal(msg.Body, q)
	if err != nil {
		client.log.Warn("is no iq bind: ", err)
		return nil, client
	}
	if q.Resource == "" {
		client.jid.Resource = makeResource()
	} else {
		client.jid.Resource = q.Resource
	}
	client.log = client.log.WithField("jid", client.jid.Full())
	client.out.Encode(&messages.IQ{
		Type: messages.IQTypeResult,
		ID:   msg.ID,
		Body: []byte(fmt.Sprintf(
			`<bind xmlns='%s'>
				<jid>%s</jid>
			</bind>`,
			messages.NSBind, client.jid.Full())),
	})

	return state.Next, client
}

// SendingClient state
type SendingClient struct {
	Next State
}

// Process messages
func (state *SendingClient) Process(client *Client) (State, *Client) {
	client.log = client.log.WithField("state", "normal")
	client.log.Debug("sending")
	// sending
	go func() {
		select {
		case msg := <-client.messages:
			err := client.out.Encode(msg)
			client.log.Info(err)
		case <-client.close:
			return
		}
	}()
	client.log.Debug("receiving")
	return state.Next, client
}

// ReceivingClient state
type ReceivingClient struct {
}

// Process messages
func (state *ReceivingClient) Process(client *Client) (State, *Client) {
	element, err := client.Read()
	if err != nil {
		client.log.Warn("unable to read: ", err)
		return nil, client
	}
	/*
		for _, extension := range client.Server.Extensions {
			extension.Process(element, client)
		}*/
	client.log.Debug(element)
	return state, client
}
