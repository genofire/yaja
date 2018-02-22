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

	Bind              *Bind                     // RFC 6120: A.7  Resource binding namespace (But in a IQ?)
	Roster            *xmppiq.RosterQuery       // RFC 6121: Appendix D.
	DiscoQueryInfo    *xmppiq.DiscoQueryInfo    // XEP-0030: XMPP Service Discovery (see iq/service_discovery.go)
	DiscoQueryItem    *xmppiq.DiscoQueryItem    // XEP-0030: XMPP Service Discovery (see iq/service_discovery.go)
	PrivateXMLStorage *xmppiq.PrivateXMLStorage // XEP-0049: Private XML Storage (see iq/private_xml_storage.go)
	VCard             *xmppiq.VCard             // XEP-0054: vcard-temp (see iq/vcard.go) - WIP
	Register          *xmppiq.Register          // XEP-0077: In-Band Registration - 14.1 (see iq/register.go) - WIP
	Version           *xmppiq.Version           // XEP-0092: Software Version (see iq/version.go)
	Ping              *xmppiq.Ping              // XEP-0199: XMPP Ping (see iq/ping.go)

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
