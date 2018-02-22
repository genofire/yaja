package xmppiq

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

// RosterQuery implements RFC 6121 - Appendix D.  XML Schema for jabber:iq:roster
type RosterQuery struct {
	XMLName xml.Name     `xml:"jabber:iq:roster query"`
	Version string       `xml:"ver,attr,omitempty"`
	Items   []RosterItem `xml:"item"`
}

// RosterAskType is a Enum of item attribute ask
type RosterAskType string

// RFC 6121: Appendix D.  XML Schema for jabber:iq:roster
const (
	RosterAskSubscribe RosterAskType = "subscribe"
	RosterAskNone      RosterAskType = ""
)

// RosterAskType is a Enum of item attribute subscription
type RosterSubscriptionType string

// RFC 6121: Appendix D.  XML Schema for jabber:iq:roster
const (
	RosterSubscriptionBoth   RosterSubscriptionType = "both"
	RosterSubscriptionFrom   RosterSubscriptionType = "from"
	RosterSubscriptionNone   RosterSubscriptionType = "none"
	RosterSubscriptionRemove RosterSubscriptionType = "remove"
	RosterSubscriptionTo     RosterSubscriptionType = "to"
)

// RosterItem implements RFC 6121 - Appendix D.  XML Schema for jabber:iq:roster
type RosterItem struct {
	JID          *xmppbase.JID          `xml:"jid,attr"`
	Approved     *bool                  `xml:"approved,attr,omitempty"`
	Ask          RosterAskType          `xml:"ask,attr,omitempty"`
	Name         string                 `xml:"name,attr,omitempty"`
	Subscription RosterSubscriptionType `xml:"subscription,attr"`
	Group        []string               `xml:"group"`
}
