package xmppiq

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

// DiscoQueryInfo implements XEP 0030: Service Discovery - 11.1 disco#info
type DiscoQueryInfo struct {
	XMLName    xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
	Node       *string  `xml:"node,attr"`
	Identities []*DiscoIdentity
	Features   []*DiscoFeature
}

// DiscoIdentity implements XEP 0030: Service Discovery - 11.1 disco#info
type DiscoIdentity struct {
	XMLName  xml.Name `xml:"http://jabber.org/protocol/disco#info identity"`
	Category string   `xml:"category"` //required
	Name     *string  `xml:"name"`
	Type     string   `xml:"type"` //required
}

// DiscoFeature implements XEP 0030: Service Discovery - 11.1 disco#info
type DiscoFeature struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info feature"`
	Var     string   `xml:"var"` //required
}

// DiscoQueryItem implements XEP 0030: Service Discovery - 11.2 disco#items
type DiscoQueryItem struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#items query"`
	Node    *string  `xml:"node,attr"`
	Items   []*DiscoItem
}

// DiscoItem implements XEP 0030: Service Discovery - 11.2 disco#items
type DiscoItem struct {
	XMLName xml.Name      `xml:"http://jabber.org/protocol/disco#items item"`
	JID     *xmppbase.JID `xml:"jid"`
	Node    *string       `xml:"node"`
	Name    *string       `xml:"name"`
}
