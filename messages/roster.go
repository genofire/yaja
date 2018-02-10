package messages

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/model"
)

type ClientQuery struct {
	Item []RosterItem
}

type RosterItem struct {
	XMLName      xml.Name   `xml:"jabber:iq:roster item"`
	JID          *model.JID `xml:",attr"`
	Name         string     `xml:",attr"`
	Subscription string     `xml:",attr"`
	Group        []string
}
