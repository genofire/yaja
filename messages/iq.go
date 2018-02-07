package messages

import "encoding/xml"

type IQType string

const (
	IQTypeGet    IQType = "get"
	IQTypeSet    IQType = "set"
	IQTypeResult IQType = "result"
	IQTypeError  IQType = "error"
)

// IQ element - info/query
type IQClient struct {
	XMLName xml.Name     `xml:"jabber:client iq"`
	From    string       `xml:"from,attr"`
	ID      string       `xml:"id,attr"`
	To      string       `xml:"to,attr"`
	Type    IQType       `xml:"type,attr"`
	Error   *ErrorClient `xml:"error"`
	Bind    Bind         `xml:"bind"`
	Body    []byte       `xml:",innerxml"`
	// RosterRequest - better detection of iq's
}
