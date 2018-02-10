package client

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
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

func NewClient(jid *model.JID, password string) (*Client, error) {
	_, srvEntries, err := net.LookupSRV("xmpp-client", "tcp", jid.Domain)
	addr := jid.Domain + ":5222"
	if err == nil && len(srvEntries) > 0 {
		bestSrv := srvEntries[0]
		for _, srv := range srvEntries {
			if srv.Priority <= bestSrv.Priority && srv.Weight >= bestSrv.Weight {
				bestSrv = srv
				addr = fmt.Sprintf("%s:%d", srv.Target, srv.Port)
			}
		}
	}
	a := strings.SplitN(addr, ":", 2)
	if len(a) == 1 {
		addr += ":5222"
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	client := &Client{
		conn: conn,
		In:   xml.NewDecoder(conn),
		Out:  xml.NewEncoder(conn),

		JID: jid,
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

func (client *Client) startStream() (*messages.StreamFeatures, error) {
	// XMPP-Connection
	_, err := fmt.Fprintf(client.conn, "<?xml version='1.0'?>\n"+
		"<stream:stream to='%s' xmlns='%s'\n"+
		" xmlns:stream='%s' version='1.0'>\n",
		model.XMLEscape(client.JID.Domain), messages.NSClient, messages.NSStream)
	if err != nil {
		return nil, err
	}
	element, err := client.Read()
	if err != nil {
		return nil, err
	}
	if element.Name.Space != messages.NSStream || element.Name.Local != "stream" {
		return nil, errors.New("is not stream")
	}
	f := &messages.StreamFeatures{}
	if err := client.ReadElement(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (client *Client) connect(password string) error {
	if _, err := client.startStream(); err != nil {
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
	err := client.auth(password)
	if err != nil {
		return err
	}
	_, err = client.startStream()
	if err != nil {
		return err
	}
	// bind to resource
	bind := &messages.Bind{}
	if client.JID.Resource != "" {
		bind.Resource = client.JID.Resource
	}
	client.Out.Encode(&messages.IQClient{
		Type: messages.IQTypeSet,
		To:   model.NewJID(client.JID.Domain),
		From: client.JID,
		ID:   utils.CreateCookieString(),
		Bind: bind,
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
		return errors.New(fmt.Sprintf("%v", iq.Other))
	}
	// set status
	client.Send(&messages.PresenceClient{Show: "online", Status: "yaja client"})

	return err
}
