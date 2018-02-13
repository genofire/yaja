package messages

// which XEP ????

import (
	"encoding/xml"
)

// where to put: (server part debug? is it part)

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
