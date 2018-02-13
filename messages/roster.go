package messages

// which XEP ????

import (
	"encoding/xml"
)

type ClientQuery struct {
	Item []RosterItem
}

type RosterItem struct {
	XMLName      xml.Name `xml:"jabber:iq:roster item"`
	JID          *JID     `xml:",attr"`
	Name         string   `xml:",attr"`
	Subscription string   `xml:",attr"`
	Group        []string
}
