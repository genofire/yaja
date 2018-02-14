package xmppiq

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

// IQDiscoQueryInfo implements XEP 0030: Service Discovery - 11.1 disco#info
type IQDiscoQueryInfo struct {
	XMLName    xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
	Node       *string  `xml:"node,attr"`
	Identities []*IQDiscoIdentity
	Features   []*IQDiscoFeature
}

// IQDiscoIdentity implements XEP 0030: Service Discovery - 11.1 disco#info
type IQDiscoIdentity struct {
	XMLName  xml.Name `xml:"http://jabber.org/protocol/disco#info identity"`
	Category string   `xml:"category"` //required
	Name     *string  `xml:"name"`
	Type     string   `xml:"type"` //required
}

// IQDiscoFeature implements XEP 0030: Service Discovery - 11.1 disco#info
type IQDiscoFeature struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info feature"`
	Var     string   `xml:"var"` //required
}

// IQDiscoQueryItem implements XEP 0030: Service Discovery - 11.2 disco#items
type IQDiscoQueryItem struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#items query"`
	Node    *string  `xml:"node,attr"`
	Items   []*IQDiscoItem
}

// IQDiscoItem implements XEP 0030: Service Discovery - 11.2 disco#items
type IQDiscoItem struct {
	XMLName xml.Name      `xml:"http://jabber.org/protocol/disco#items item"`
	JID     *xmppbase.JID `xml:"jid"`
	Node    *string       `xml:"node"`
	Name    *string       `xml:"name"`
}
