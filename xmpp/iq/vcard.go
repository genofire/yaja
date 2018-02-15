package xmppiq

import (
	"encoding/xml"
)

// WARNING WIP

// VCard implements XEP-0054: vcard-temp - 14
//TODO
type VCard struct {
	XMLName xml.Name `xml:"vcard-temp vCard"`
	Body    []byte   `xml:",innerxml"`
}
