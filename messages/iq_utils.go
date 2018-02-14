package messages

import (
	"encoding/xml"
)

// XEP-0199: XMPP Ping
type Ping struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}

// which XEP ????
// where to put: (server part debug? is it part)
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
