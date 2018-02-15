package xmppiq

import (
	"encoding/xml"
)

// Ping implements XEP-0199: XMPP Ping - 10
type Ping struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}
