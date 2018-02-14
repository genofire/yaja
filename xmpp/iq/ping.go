package xmppiq

import (
	"encoding/xml"
)

// XEP-0199: XMPP Ping
type Ping struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}
