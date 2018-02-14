package xmpp

import "encoding/xml"

// ErrorServer implements RFC 6120 - A.6 Server Namespace (a part)
type ErrorServer struct {
	XMLName xml.Name  `xml:"jabber:server error"`
	Code    string    `xml:"code,attr,omitempty"`
	Type    ErrorType `xml:"type,attr"` // required
	Text    *Text

	// RFC 6120  A.8  Resource binding namespace
	StanzaErrorGroup

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
