package xmppiq

import (
	"encoding/xml"
)

// Version implements XEP-0092: Software Version - 4
type Version struct {
	XMLName xml.Name `xml:"jabber:iq:version query"`
	Name    string   `xml:"name"`    //required
	Version string   `xml:"version"` //required
	OS      *string  `xml:"os"`
}
