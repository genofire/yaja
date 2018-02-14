package xmpp

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

// PresenceClient implements RFC 6120 - A.5 Client Namespace (a part)
type PresenceClient struct {
	XMLName xml.Name      `xml:"jabber:client presence"`
	From    *xmppbase.JID `xml:"from,attr,omitempty"`
	ID      string        `xml:"id,attr,omitempty"`
	To      *xmppbase.JID `xml:"to,attr,omitempty"`
	Type    PresenceType  `xml:"type,attr,omitempty"`
	Lang    string        `xml:"lang,attr,omitempty"`

	Show     PresenceShow `xml:"show,omitempty"`
	Status   string       `xml:"status,omitempty"`
	Priority uint         `xml:"priority,omitempty"` // default: 0

	Error *ErrorClient

	Delay *Delay `xml:"delay"` // which XEP ?

	// which XEP ?
	// Caps     *ClientCaps  `xml:"c"`

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
