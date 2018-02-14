package xmpp

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

// ClientQuery implements which XEP ????
type ClientQuery struct {
	Item []RosterItem
}

// RosterItem implements which XEP ????
type RosterItem struct {
	XMLName      xml.Name      `xml:"jabber:iq:roster item"`
	JID          *xmppbase.JID `xml:",attr"`
	Name         string        `xml:",attr"`
	Subscription string        `xml:",attr"`
	Group        []string
}
