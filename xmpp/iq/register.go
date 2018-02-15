package xmppiq

import (
	"encoding/xml"
)

// WARNING WIP

// Register implements XEP-0077: In-Band Registration - 14.1 jabber:iq:register
//TODO
type Register struct {
	XMLName      xml.Name `xml:"jabber:iq:register query"`
	Instructions string   `xml:"instructions"`
	Username     string   `xml:"username"`
	Password     string   `xml:"password"`
}

// FeatureRegister implements XEP-0077: In-Band Registration - 14.2 Stream Feature
type FeatureRegister struct {
	XMLName xml.Name `xml:"http://jabber.org/features/iq-register register"`
}
