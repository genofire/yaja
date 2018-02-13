package messages

import (
	"encoding/xml"
)

// RFC 6120 - A.5 Client Namespace (a part)
type MessageClient struct {
	XMLName xml.Name    `xml:"jabber:client message"`
	From    *JID        `xml:"from,attr,omitempty"`
	ID      string      `xml:"id,attr,omitempty"`
	To      *JID        `xml:"to,attr,omitempty"`
	Type    MessageType `xml:"type,attr,omitempty"` // default: normal
	Lang    string      `xml:"lang,attr,omitempty"`

	Subject string `xml:"subject,omitempty"`
	Body    string `xml:"body,omitempty"`
	Thread  string `xml:"thread,omitempty"`
	Error   *ErrorClient

	Delay *Delay `xml:"delay"` // which XEP ?

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
