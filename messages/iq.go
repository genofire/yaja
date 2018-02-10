package messages

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/model"
)

type IQType string

const (
	IQTypeGet    IQType = "get"
	IQTypeSet    IQType = "set"
	IQTypeResult IQType = "result"
	IQTypeError  IQType = "error"
)

// IQ element - info/query
type IQClient struct {
	XMLName         xml.Name     `xml:"jabber:client iq"`
	From            *model.JID   `xml:"from,attr"`
	ID              string       `xml:"id,attr"`
	To              *model.JID   `xml:"to,attr"`
	Type            IQType       `xml:"type,attr"`
	Error           *ErrorClient `xml:"error"`
	Bind            *Bind
	Ping            *Ping
	PrivateQuery    *IQPrivateQuery
	PrivateRegister *IQPrivateRegister
	// Any hasn't matched element
	Other []XMLElement `xml:",any"`
}

type Ping struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}

type IQPrivateQuery struct {
	XMLName xml.Name `xml:"jabber:iq:private query"`
	Body    []byte   `xml:",innerxml"`
}

type IQPrivateRegister struct {
	XMLName      xml.Name `xml:"jabber:iq:register query"`
	Instructions string   `xml:"instructions"`
	Username     string   `xml:"username"`
	Password     string   `xml:"password"`
}
