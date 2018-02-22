package xmpp

import "encoding/xml"

// ErrorClient implements RFC 6120 - A.5 Client Namespace (a part)
type ErrorClient struct {
	XMLName xml.Name  `xml:"jabber:client error"`
	Code    string    `xml:"code,attr,omitempty"`
	Type    ErrorType `xml:"type,attr"` // required
	Text    *Text

	StanzaErrorGroup // RFC 6120: A.8  Resource binding namespace

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
