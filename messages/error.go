package messages

import "encoding/xml"

type StreamError struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams error"`
	Text    string

	Any xml.Name `xml:",any"`
}

// ErrorClient element
type ErrorClient struct {
	XMLName xml.Name `xml:"jabber:client error"`
	Code    string   `xml:"code,attr"`
	Type    string   `xml:"type,attr"`
	Text    string   `xml:"text"`

	Any xml.Name `xml:",any"`
}
