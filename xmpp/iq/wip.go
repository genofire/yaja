package xmppiq

import (
	"encoding/xml"
)

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
