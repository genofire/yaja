package xmpp

import "encoding/xml"

// Text implements RFC 6120 part of A.2 Stream Error Namespace, A.4 SASL Namespace and A.8 Stanza Error Namespace
type Text struct {
	// `xml:"urn:ietf:params:xml:ns:xmpp-streams text"`
	// `xml:"urn:ietf:params:xml:ns:xmpp-sasl text"`
	// `xml:"urn:ietf:params:xml:ns:xmpp-stanzas text"`
	XMLName xml.Name `xml:"text"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Body    string   `xml:",innerxml"`
}

// XMLElement is for Unmarshal undefined structs a fallback - any hasn't matched element
type XMLElement struct {
	XMLName  xml.Name
	InnerXML string `xml:",innerxml"`
}

// Delay implements which XEP ?
type Delay struct {
	Stamp string `xml:"stamp,attr"`
}
