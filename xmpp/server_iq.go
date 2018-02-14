package xmpp

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
	"dev.sum7.eu/genofire/yaja/xmpp/iq"
)

// IQServer implements RFC 6120 - A.6 Server Namespace (a part)
type IQServer struct {
	XMLName xml.Name      `xml:"jabber:server iq"`
	From    *xmppbase.JID `xml:"from,attr,omitempty"`
	ID      string        `xml:"id,attr"`   // required
	To      *xmppbase.JID `xml:"to,attr"`   // required
	Type    IQType        `xml:"type,attr"` // required
	Error   *ErrorServer

	Bind *Bind        // RFC 6120  A.7  Resource binding namespace (But in a IQ?)
	Ping *xmppiq.Ping // XEP-0199: XMPP Ping (see iq/ping.go)

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
