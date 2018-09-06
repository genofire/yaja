package xmuc

import (
	"encoding/xml"
	"time"
)

// Base implements XEP-0045: Multi-User Chat - 19.1
type Base struct {
	XMLName  xml.Name `xml:"http://jabber.org/protocol/muc x"`
	History  *History `xml:"history,omitempty"`
	Password string   `xml:"password,omitempty"`
}

// History implements XEP-0045: Multi-User Chat - 19.1
type History struct {
	MaxChars   *int       `xml:"maxchars,attr,omitempty"`
	MaxStanzas *int       `xml:"maxstanzas,attr,omitempty"`
	Seconds    *int       `xml:"seconds,attr,omitempty"`
	Since      *time.Time `xml:"since,attr,omitempty"`
}
