package messages

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/model"
)

type XMLElement struct {
	XMLName  xml.Name
	InnerXML string `xml:",innerxml"`
}

type Delay struct {
	Stamp string `xml:"stamp,attr"`
}

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

type ShowType string

const (
	ShowTypeAway ShowType = "away"
	ShowTypeChat ShowType = "chat"
	ShowTypeDND  ShowType = "dnd"
	ShowTypeXA   ShowType = "xa"
)

// PresenceClient element
type PresenceClient struct {
	XMLName xml.Name     `xml:"jabber:client presence"`
	From    *model.JID   `xml:"from,attr,omitempty"`
	ID      string       `xml:"id,attr,omitempty"`
	To      *model.JID   `xml:"to,attr,omitempty"`
	Type    PresenceType `xml:"type,attr,omitempty"`
	Lang    string       `xml:"lang,attr,omitempty"`

	Show     ShowType `xml:"show,omitempty"`   // away, chat, dnd, xa
	Status   string   `xml:"status,omitempty"` // sb []clientText
	Priority string   `xml:"priority,omitempty"`
	// Caps     *ClientCaps  `xml:"c"`
	Delay *Delay `xml:"delay"`

	Error *ErrorClient
}

type ChatType string

const (
	ChatTypeChat      ChatType = "chat"
	ChatTypeGroupchat ChatType = "groupchat"
	ChatTypeError     ChatType = "error"
	ChatTypeHeadline  ChatType = "headline"
	ChatTypeNormal    ChatType = "normal"
)

// MessageClient element
type MessageClient struct {
	XMLName xml.Name   `xml:"jabber:client message"`
	From    *model.JID `xml:"from,attr,omitempty"`
	ID      string     `xml:"id,attr,omitempty"`
	To      *model.JID `xml:"to,attr,omitempty"`
	Type    ChatType   `xml:"type,attr,omitempty"`
	Lang    string     `xml:"lang,attr,omitempty"`
	Subject string     `xml:"subject"`
	Body    string     `xml:"body"`
	Thread  string     `xml:"thread"`
	// Any hasn't matched element
	Other []XMLElement `xml:",any"`

	Delay *Delay `xml:"delay"`
	Error *ErrorClient
}
