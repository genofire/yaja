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
type IQ struct {
	XMLName xml.Name `xml:"jabber:client iq"`
	From    string   `xml:"from,attr"`
	ID      string   `xml:"id,attr"`
	To      string   `xml:"to,attr"`
	Type    IQType   `xml:"type,attr"`
	Error   *Error   `xml:"error"`
	//Bind    bindBind    `xml:"bind"`
	Body []byte `xml:",innerxml"`
	// RosterRequest - better detection of iq's
}
