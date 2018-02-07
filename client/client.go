package client

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strings"

	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/model"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

// Client holds XMPP connection opitons
type Client struct {
	conn net.Conn // connection to server
	Out  *xml.Encoder
	In   *xml.Decoder

	JID *model.JID
}

func NewClient(jid model.JID, password string) (*Client, error) {
	conn, err := net.Dial("tcp", jid.Domain+":5222")
	if err != nil {
		return nil, err
	}
	client := &Client{
		conn: conn,
		In:   xml.NewDecoder(conn),
		Out:  xml.NewEncoder(conn),

		JID: &jid,
	}

	if err = client.connect(password); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

// Close closes the XMPP connection
func (c *Client) Close() error {
	if c.conn != (*tls.Conn)(nil) {
		return c.conn.Close()
	}
	return nil
}

func (client *Client) Read() (*xml.StartElement, error) {
	for {
		nextToken, err := client.In.Token()
		if err != nil {
			return nil, err
		}
		switch nextToken.(type) {
		case xml.StartElement:
			element := nextToken.(xml.StartElement)
			return &element, nil
		}
	}
}
func (client *Client) ReadElement(p interface{}) error {
	element, err := client.Read()
	if err != nil {
		return err
	}
	return client.In.DecodeElement(p, element)
}

func (client *Client) init() error {
	// XMPP-Connection
	_, err := fmt.Fprintf(client.conn, "<?xml version='1.0'?>\n"+
		"<stream:stream to='%s' xmlns='%s'\n"+
		" xmlns:stream='%s' version='1.0'>\n",
		model.XMLEscape(client.JID.Domain), messages.NSClient, messages.NSStream)
	if err != nil {
		return err
	}
	element, err := client.Read()
	if err != nil {
		return err
	}
	if element.Name.Space != messages.NSStream || element.Name.Local != "stream" {
		return errors.New("is not stream")
	}
	return nil
}
func (client *Client) connect(password string) error {
	if err := client.init(); err != nil {
		return err
	}
	var f messages.StreamFeatures
	if err := client.ReadElement(&f); err != nil {
		return err
	}
	if err := client.Out.Encode(&messages.TLSStartTLS{}); err != nil {
		return err
	}

	var p messages.TLSProceed
	if err := client.ReadElement(&p); err != nil {
		return err
	}
	// Change tcp to tls
	tlsconn := tls.Client(client.conn, &tls.Config{
		ServerName: client.JID.Domain,
	})
	client.conn = tlsconn
	client.In = xml.NewDecoder(client.conn)
	client.Out = xml.NewEncoder(client.conn)

	if err := tlsconn.Handshake(); err != nil {
		return err
	}
	if err := tlsconn.VerifyHostname(client.JID.Domain); err != nil {
		return err
	}
	if err := client.init(); err != nil {
		return err
	}
	//auth:
	if err := client.ReadElement(&f); err != nil {
		return err
	}
	mechanism := ""
	for _, m := range f.Mechanisms.Mechanism {
		if m == "PLAIN" {
			mechanism = m
			// Plain authentication: send base64-encoded \x00 user \x00 password.
			raw := "\x00" + client.JID.Local + "\x00" + password
			enc := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
			base64.StdEncoding.Encode(enc, []byte(raw))
			fmt.Fprintf(client.conn, "<auth xmlns='%s' mechanism='PLAIN'>%s</auth>\n", messages.NSSASL, enc)
			break
		}
		if m == "DIGEST-MD5" {
			mechanism = m
			// Digest-MD5 authentication
			fmt.Fprintf(client.conn, "<auth xmlns='%s' mechanism='DIGEST-MD5'/>\n", messages.NSSASL)
			var ch string
			if err := client.ReadElement(&ch); err != nil {
				return err
			}
			b, err := base64.StdEncoding.DecodeString(string(ch))
			if err != nil {
				return err
			}
			tokens := map[string]string{}
			for _, token := range strings.Split(string(b), ",") {
				kv := strings.SplitN(strings.TrimSpace(token), "=", 2)
				if len(kv) == 2 {
					if kv[1][0] == '"' && kv[1][len(kv[1])-1] == '"' {
						kv[1] = kv[1][1 : len(kv[1])-1]
					}
					tokens[kv[0]] = kv[1]
				}
			}
			realm, _ := tokens["realm"]
			nonce, _ := tokens["nonce"]
			qop, _ := tokens["qop"]
			charset, _ := tokens["charset"]
			cnonceStr := cnonce()
			digestURI := "xmpp/" + client.JID.Domain
			nonceCount := fmt.Sprintf("%08x", 1)
			digest := saslDigestResponse(client.JID.Local, realm, password, nonce, cnonceStr, "AUTHENTICATE", digestURI, nonceCount)
			message := "username=\"" + client.JID.Local + "\", realm=\"" + realm + "\", nonce=\"" + nonce + "\", cnonce=\"" + cnonceStr +
				"\", nc=" + nonceCount + ", qop=" + qop + ", digest-uri=\"" + digestURI + "\", response=" + digest + ", charset=" + charset

			fmt.Fprintf(client.conn, "<response xmlns='%s'>%s</response>\n", messages.NSSASL, base64.StdEncoding.EncodeToString([]byte(message)))

			err = client.ReadElement(&ch)
			if err != nil {
				return err
			}
			_, err = base64.StdEncoding.DecodeString(ch)
			if err != nil {
				return err
			}
			fmt.Fprintf(client.conn, "<response xmlns='%s'/>\n", messages.NSSASL)
			break
		}
	}
	if mechanism == "" {
		return fmt.Errorf("PLAIN authentication is not an option: %v", f.Mechanisms.Mechanism)
	}
	element, err := client.Read()
	if err != nil {
		return err
	}
	if element.Name.Local != "success" {
		return errors.New("auth failed: " + element.Name.Local)
	}

	err = client.init()
	if err != nil {
		return err
	}
	if err := client.ReadElement(&f); err != nil {
		return err
	}
	// bind to resource
	var msg string
	if client.JID.Resource == "" {
		msg = fmt.Sprintf("<bind xmlns='%s'></bind>", messages.NSBind)
	} else {
		msg = fmt.Sprintf(
			`<bind xmlns='%s'>
				<resource>%s</resource>
			</bind>`,
			messages.NSBind, client.JID.Resource)
	}
	client.Out.Encode(&messages.IQClient{
		Type: messages.IQTypeSet,
		To:   client.JID.Domain,
		From: client.JID.Full(),
		ID:   utils.CreateCookieString(),
		Body: []byte(msg),
	})

	var iq messages.IQClient
	if err := client.ReadElement(&iq); err != nil {
		return err
	}
	if &iq.Bind == nil {
		return errors.New("<iq> result missing <bind>")
	}
	if iq.Bind.JID != nil {
		client.JID.Local = iq.Bind.JID.Local
		client.JID.Domain = iq.Bind.JID.Domain
		client.JID.Resource = iq.Bind.JID.Resource
	} else {
		return errors.New(string(iq.Body))
	}
	// set status
	client.Out.Encode(&messages.PresenceClient{Show: "online", Status: "yaja client"})

	return nil
}

func saslDigestResponse(username, realm, passwd, nonce, cnonceStr, authenticate, digestURI, nonceCountStr string) string {
	h := func(text string) []byte {
		h := md5.New()
		h.Write([]byte(text))
		return h.Sum(nil)
	}
	hex := func(bytes []byte) string {
		return fmt.Sprintf("%x", bytes)
	}
	kd := func(secret, data string) []byte {
		return h(secret + ":" + data)
	}

	a1 := string(h(username+":"+realm+":"+passwd)) + ":" + nonce + ":" + cnonceStr
	a2 := authenticate + ":" + digestURI
	response := hex(kd(hex(h(a1)), nonce+":"+nonceCountStr+":"+cnonceStr+":auth:"+hex(h(a2))))
	return response
}

func cnonce() string {
	randSize := big.NewInt(0)
	randSize.Lsh(big.NewInt(1), 64)
	cn, err := rand.Int(rand.Reader, randSize)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%016x", cn)
}
