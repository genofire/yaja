package xmpp

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
	"dev.sum7.eu/genofire/yaja/xmpp/iq"
)

// IQClient implements RFC 6120 - A.5 Client Namespace (a part)
type IQClient struct {
	XMLName xml.Name      `xml:"jabber:client iq"`
	From    *xmppbase.JID `xml:"from,attr,omitempty"`
	ID      string        `xml:"id,attr"` // required
	To      *xmppbase.JID `xml:"to,attr,omitempty"`
	Type    IQType        `xml:"type,attr"` // required
	Error   *ErrorClient

	Bind            *Bind                     // RFC 6120  A.7  Resource binding namespace (But in a IQ?)
	Ping            *xmppiq.Ping              // XEP-0199: XMPP Ping (see iq/ping.go)
	PrivateQuery    *xmppiq.IQPrivateQuery    // which XEP ?
	PrivateRegister *xmppiq.IQPrivateRegister // which XEP ?

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
