package messages

import (
	"encoding/xml"
)

// RFC 6120 - A.5 Client Namespace (a part)
type IQClient struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	From    *JID     `xml:"from,attr,omitempty"`
	ID      string   `xml:"id,attr"` // required
	To      *JID     `xml:"to,attr,omitempty"`
	Type    IQType   `xml:"type,attr"` // required
	Error   *ErrorClient

	Bind            *Bind              // which XEP ?
	Ping            *Ping              // which XEP ?
	PrivateQuery    *IQPrivateQuery    // which XEP ?
	PrivateRegister *IQPrivateRegister // which XEP ?

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
