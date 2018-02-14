package messages

import "encoding/xml"

// RFC 6120 - A.5 Client Namespace (a part)
type ErrorClient struct {
	XMLName xml.Name  `xml:"jabber:client error"`
	Code    string    `xml:"code,attr,omitempty"`
	Type    ErrorType `xml:"type,attr"` // required
	Text    *Text

	// RFC 6120  A.8  Resource binding namespace
	StanzaErrorGroup

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}