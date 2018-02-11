package messages

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/model"
)

// RFC 3920  C.1  Streams name space
type StreamFeatures struct {
	XMLName    xml.Name `xml:"http://etherx.jabber.org/streams features"`
	StartTLS   *TLSStartTLS
	Mechanisms SASLMechanisms
	Bind       *Bind
	Session    bool
}

// RFC 3920  C.3  TLS name space
type TLSStartTLS struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required *string  `xml:"required"`
}

type TLSFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}

type TLSProceed struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

// RFC 3920  C.5  Resource binding name space
type Bind struct {
	XMLName  xml.Name   `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string     `xml:"resource"`
	JID      *model.JID `xml:"jid"`
}
