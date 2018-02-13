package messages

import (
	"encoding/xml"
)

// RFC 6120 - A.1 Stream Namespace
type StreamFeatures struct {
	XMLName    xml.Name `xml:"http://etherx.jabber.org/streams features"`
	StartTLS   *TLSStartTLS
	Mechanisms SASLMechanisms
	Bind       *Bind
	Session    bool
}

// RFC 6120 - A.3 StartTLS Namespace
type TLSStartTLS struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required *string  `xml:"required"`
}

// RFC 6120 - A.3 StartTLS Namespace
type TLSProceed struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

// RFC 6120 - A.3 StartTLS Namespace
type TLSFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}

// RFC 6120  A.7  Resource binding namespace
type Bind struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string   `xml:"resource"`
	JID      *JID     `xml:"jid"`
}
