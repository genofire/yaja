package xmppiq

import (
	"encoding/xml"
)

// Version implements XEP-0092: Software Version - 4
type Version struct {
	XMLName xml.Name `xml:"jabber:iq:version query"`
	Name    string   `xml:"name,omitempty"`    //required
	Version string   `xml:"version,omitempty"` //required
	OS      string   `xml:"os,omitempty"`
}
