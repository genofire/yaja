package messages

import (
	"encoding/xml"
)

// RFC 6120 - A.6 Server Namespace (a part)
type IQServer struct {
	XMLName xml.Name `xml:"jabber:server iq"`
	From    *JID     `xml:"from,attr,omitempty"`
	ID      string   `xml:"id,attr"`   // required
	To      *JID     `xml:"to,attr"`   // required
	Type    IQType   `xml:"type,attr"` // required
	Error   *ErrorServer

	Bind            *Bind              // which XEP ?
	Ping            *Ping              // which XEP ?
	PrivateQuery    *IQPrivateQuery    // which XEP ?
	PrivateRegister *IQPrivateRegister // which XEP ?

	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}
