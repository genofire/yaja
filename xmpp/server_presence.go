package xmpp

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

// PresenceServer implements RFC 6120 - A.6 Server Namespace (a part)
type PresenceServer struct {
	XMLName xml.Name      `xml:"jabber:server presence"`
	From    *xmppbase.JID `xml:"from,attr"` // required
	ID      string        `xml:"id,attr,omitempty"`
	To      *xmppbase.JID `xml:"to,attr"` // required
	Type    PresenceType  `xml:"type,attr,omitempty"`
	Lang    string        `xml:"lang,attr,omitempty"`

	Show     PresenceShow `xml:"show,omitempty"`
	Status   string       `xml:"status,omitempty"`
	Priority uint         `xml:"priority,omitempty"` // default: 0

	Error *ErrorServer

	Delay *Delay `xml:"delay"` // which XEP ?

	// which XEP ?
	// Caps     *ClientCaps  `xml:"c"`

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
