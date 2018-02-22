package client

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"net"

	"dev.sum7.eu/genofire/yaja/model"
	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

func (client *Client) setConnection(conn net.Conn) {
	client.conn = conn
	client.in = xml.NewDecoder(client.conn)
	client.out = xml.NewEncoder(client.conn)
}

func (client *Client) startStream() (*xmpp.StreamFeatures, error) {
	logCTX := client.Logging.WithField("type", "stream")
	// XMPP-Connection
	_, err := fmt.Fprintf(client.conn, "<?xml version='1.0'?>\n"+
		"<stream:stream to='%s' xmlns='%s'\n"+
		" xmlns:stream='%s' version='1.0'>\n",
		model.XMLEscape(client.JID.Domain), xmpp.NSClient, xmpp.NSStream)
	if err != nil {
		return nil, err
	}
	element, err := client.Read()
	if err != nil {
		return nil, err
	}
	if element.Name.Space != xmpp.NSStream || element.Name.Local != "stream" {
		return nil, errors.New("is not stream")
	}
	f := &xmpp.StreamFeatures{}
	if err := client.ReadDecode(f); err != nil {
		return nil, err
	}
	debug := "start >"
	if f.StartTLS != nil {
		debug += " tls"
	}
	debug += " mechanism("
	mFirst := true
	for _, m := range f.Mechanisms.Mechanism {
		if mFirst {
			mFirst = false
			debug += m
		} else {
			debug += ", " + m
		}
	}
	debug += ")"
	if f.Bind != nil {
		debug += " bind"
	}
	logCTX.Info(debug)
	return f, nil
}

func (client *Client) connect(password string) error {
	if _, err := client.startStream(); err != nil {
		return err
	}
	if err := client.Send(&xmpp.TLSStartTLS{}); err != nil {
		return err
	}

	var p xmpp.TLSProceed
	if err := client.ReadDecode(&p); err != nil {
		return err
	}
	// Change tcp to tls
	tlsconn := tls.Client(client.conn, &tls.Config{
		ServerName: client.JID.Domain,
	})
	client.setConnection(tlsconn)

	if err := tlsconn.Handshake(); err != nil {
		return err
	}
	if err := tlsconn.VerifyHostname(client.JID.Domain); err != nil {
		return err
	}
	if err := client.auth(password); err != nil {
		return err
	}

	f, err := client.startStream()
	if err != nil {
		return err
	}
	bind := f.Bind
	if f.Bind == nil || (f.Bind.JID == nil && f.Bind.Resource == "") {
		// bind to resource
		if client.JID.Resource != "" {
			bind.Resource = client.JID.Resource
		}
		if err := client.Send(&xmpp.IQClient{
			Type: xmpp.IQTypeSet,
			To:   xmppbase.NewJID(client.JID.Domain),
			Bind: bind,
		}); err != nil {
			return err
		}

		var iq xmpp.IQClient
		if err := client.ReadDecode(&iq); err != nil {
			return err
		}
		if iq.Error != nil {
			return errors.New(fmt.Sprintf("recv error on iq>bind: %s[%s]: %s -> %s -> %s", iq.Error.Code, iq.Error.Type, iq.Error.Text, xmpp.XMLChildrenString(iq.Error.StanzaErrorGroup), xmpp.XMLChildrenString(iq.Error.Other)))
		} else if iq.Bind == nil {
			return errors.New("iq>bind is nil :" + xmpp.XMLChildrenString(iq.Other))
		}
		bind = iq.Bind
	}
	if bind == nil {
		return errors.New("bind is nil")
	} else if bind.JID != nil {
		client.JID.Local = bind.JID.Local
		client.JID.Domain = bind.JID.Domain
		client.JID.Resource = bind.JID.Resource
		client.Logging.WithField("type", "bind").Infof("set jid by server bind '%s'", bind.JID.Full())
	} else if bind.Resource != "" {
		client.JID.Resource = bind.Resource
		client.Logging.WithField("type", "bind").Infof("set resource by server bind '%s'", bind.Resource)
	} else {
		return errors.New("bind>jid is nil" + xmpp.XMLChildrenString(bind))
	}
	// set status
	return client.Send(&xmpp.PresenceClient{Show: xmpp.PresenceShowXA, Status: "online"})
}
