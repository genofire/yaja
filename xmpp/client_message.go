package xmpp

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

// MessageClient implements RFC 6120 - A.5 Client Namespace (a part)
type MessageClient struct {
	XMLName xml.Name      `xml:"jabber:client message"`
	From    *xmppbase.JID `xml:"from,attr,omitempty"`
	ID      string        `xml:"id,attr,omitempty"`
	To      *xmppbase.JID `xml:"to,attr,omitempty"`
	Type    MessageType   `xml:"type,attr,omitempty"` // default: normal
	Lang    string        `xml:"lang,attr,omitempty"`

	Subject string `xml:"subject,omitempty"`
	Body    string `xml:"body,omitempty"`
	Thread  string `xml:"thread,omitempty"`
	Error   *ErrorClient

	Delay *Delay `xml:"delay"` // which XEP ?

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
