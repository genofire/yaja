package xmppiq

import "encoding/xml"

// PrivateXMLStorage implements XEP-0049: Private XML Storage - 7
type PrivateXMLStorage struct {
	XMLName xml.Name `xml:"jabber:iq:private query"`
	Body    []byte   `xml:",innerxml"`
}
