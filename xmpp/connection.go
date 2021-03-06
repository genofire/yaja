package xmpp

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

// StreamFeatures implements RFC 6120 - A.1 Stream Namespace
type StreamFeatures struct {
	XMLName     xml.Name `xml:"http://etherx.jabber.org/streams features"`
	StartTLS    *TLSStartTLS
	Mechanisms  SASLMechanisms      // RFC 6120: A.4 SASL Namespace (see sasl.go)
	Compression *CompressionFeature // XEP-0138: Stream Compression (see compression.go)
	Bind        *Bind
	Session     bool
}

// TLSStartTLS implements RFC 6120 - A.3 StartTLS Namespace
type TLSStartTLS struct {
	XMLName  xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
	Required string   `xml:"required,omitempty"`
}

// TLSProceed implements RFC 6120 - A.3 StartTLS Namespace
type TLSProceed struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls proceed"`
}

// TLSFailure implements RFC 6120 - A.3 StartTLS Namespace
type TLSFailure struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls failure"`
}

// Bind implements RFC 6120 - A.7  Resource binding namespace
type Bind struct {
	XMLName  xml.Name      `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
	Resource string        `xml:"resource"`
	JID      *xmppbase.JID `xml:"jid"`
}
