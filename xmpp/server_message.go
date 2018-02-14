package xmpp

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

// MessageServer implements RFC 6120 - A.6 Server Namespace (a part)
type MessageServer struct {
	XMLName xml.Name      `xml:"jabber:server message"`
	From    *xmppbase.JID `xml:"from,attr"` // required
	ID      string        `xml:"id,attr,omitempty"`
	To      *xmppbase.JID `xml:"to,attr"`             // required
	Type    MessageType   `xml:"type,attr,omitempty"` // default: normal
	Lang    string        `xml:"lang,attr,omitempty"`

	Subject string `xml:"subject,omitempty"`
	Body    string `xml:"body,omitempty"`
	Thread  string `xml:"thread,omitempty"`
	Error   *ErrorServer

	Delay *Delay `xml:"delay"` // which XEP ?

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
