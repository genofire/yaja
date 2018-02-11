package messages

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/model"
)

type PresenceType string

const (
	PresenceTypeUnavailable  PresenceType = "unavailable"
	PresenceTypeSubscribe    PresenceType = "subscribe"
	PresenceTypeSubscribed   PresenceType = "subscribed"
	PresenceTypeUnsubscribe  PresenceType = "unsubscribe"
	PresenceTypeUnsubscribed PresenceType = "unsubscribed"
	PresenceTypeProbe        PresenceType = "probe"
	PresenceTypeError        PresenceType = "error"
)

type PresenceShow string

const (
	PresenceShowAway PresenceShow = "away"
	PresenceShowChat PresenceShow = "chat"
	PresenceShowDND  PresenceShow = "dnd"
	PresenceShowXA   PresenceShow = "xa"
)

// PresenceClient element
type PresenceClient struct {
	XMLName xml.Name     `xml:"jabber:client presence"`
	From    *model.JID   `xml:"from,attr,omitempty"`
	ID      string       `xml:"id,attr,omitempty"`
	To      *model.JID   `xml:"to,attr,omitempty"`
	Type    PresenceType `xml:"type,attr,omitempty"`
	Lang    string       `xml:"lang,attr,omitempty"`

	Show     PresenceShow `xml:"show,omitempty"`   // away, chat, dnd, xa
	Status   string       `xml:"status,omitempty"` // sb []clientText
	Priority string       `xml:"priority,omitempty"`
	// Caps     *ClientCaps  `xml:"c"`
	Delay *Delay `xml:"delay"`

	Error *ErrorClient
}
